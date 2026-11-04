package services

import (
	"elcomware/app/actions"
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"time"
)

type OrderService struct {
	orderQueries   *queries.OrderQueries
	checkoutAction *actions.ExecuteCheckoutAction
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderQueries:   queries.NewOrderQueries(),
		checkoutAction: actions.NewExecuteCheckoutAction(),
	}
}

func (s *OrderService) PlaceOrder(req *requests.CheckoutRequest, user *models.User) (*models.Order, error) {
	return s.checkoutAction.ExecuteFromRequest(req, user)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderQueries.FindAll()
}

func (s *OrderService) GetUserOrders(userID string) ([]models.Order, error) {
	return s.orderQueries.FindByUserID(userID)
}

// CancelOrder retains the order and its items while releasing its reservation once.
func CancelOrder(id string, actor *models.User) error {
	if actor == nil || actor.ID == "" {
		return errors.New("sign in required")
	}
	return facades.Orm().Transaction(func(tx orm.Query) error {
		var order models.Order
		if err := tx.Where("id = ?", id).First(&order); err != nil {
			return err
		}
		if order.ID == "" || !actor.Can("orders.cancel", order.UserID) {
			return errors.New("order not found or access denied")
		}
		result, err := tx.Model(&models.Order{}).Where("id = ? AND status = ?", id, models.OrderStatusPending).Update(map[string]any{"status": models.OrderStatusCancelled, "updated_at": time.Now()})
		if err != nil {
			return err
		}
		if result.RowsAffected != 1 {
			return errors.New("only pending orders can be cancelled")
		}
		var items []models.OrderItem
		if err := tx.Where("order_id = ?", id).Get(&items); err != nil {
			return err
		}
		for _, item := range items {
			if !item.StockReserved {
				continue
			}
			if item.Quantity < 1 || item.Quantity > 1000000 {
				return errors.New("order inventory could not be verified")
			}
			if _, err := tx.Exec("UPDATE products SET stock = stock + ?, updated_at = ? WHERE id = ? AND stock >= 0", item.Quantity, time.Now(), item.ProductID); err != nil {
				return err
			}
		}
		return nil
	})
}
