package actions

import (
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/database/orm"
	"math"
	"sort"
	"time"
)

type ExecuteCheckoutAction struct{}

func NewExecuteCheckoutAction() *ExecuteCheckoutAction { return &ExecuteCheckoutAction{} }
func (a *ExecuteCheckoutAction) ExecuteFromRequest(req *requests.CheckoutRequest, user *models.User) (*models.Order, error) {
	return a.ExecuteCart(map[string]int{req.ProductID: req.Quantity}, user)
}
func (a *ExecuteCheckoutAction) ExecuteCart(cart map[string]int, user *models.User) (*models.Order, error) {
	if user == nil || user.ID == "" {
		return nil, errors.New("sign in before placing an order")
	}
	if len(cart) == 0 || len(cart) > 100 {
		return nil, errors.New("cart must contain between 1 and 100 listings")
	}
	now := time.Now()
	id := uuid.NewString()
	order := &models.Order{ID: "ord_" + id, OrderNumber: "SOKO-" + id, UserID: user.ID, CustomerName: user.Name, CustomerEmail: user.Email, Currency: "USD", Status: models.OrderStatusPending, PaymentMethod: "Payment not collected", CreatedAt: now, UpdatedAt: now, Items: []models.OrderItem{}}
	ids := make([]string, 0, len(cart))
	for id := range cart {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		for _, id := range ids {
			var p models.Product
			if err := tx.Where("id = ?", id).First(&p); err != nil {
				return err
			}
			if p.ID == "" {
				return errors.New("a cart listing is no longer available")
			}
			qty := cart[id]
			price, err := p.PriceFor(qty)
			if err != nil {
				return fmt.Errorf("%s: %w", p.Title, err)
			}
			if p.Stock >= 0 {
				result, err := tx.Exec("UPDATE products SET stock = stock - ?, updated_at = ? WHERE id = ? AND stock >= ?", qty, now, id, qty)
				if err != nil {
					return err
				}
				if result.RowsAffected != 1 {
					return errors.New("stock changed; review your cart and try again")
				}
			}
			total := math.Round(price*float64(qty)*100) / 100
			if math.IsNaN(total) || math.IsInf(total, 0) || total > 1e12 || order.TotalAmount+total > 1e12 {
				return errors.New("order amount exceeds the supported limit")
			}
			item := models.OrderItem{StockReserved: p.Stock >= 0, ID: "item_" + uuid.NewString(), OrderID: order.ID, ProductID: id, ProductTitle: p.Title, ProductType: p.Type, Quantity: qty, UnitPrice: price, TotalPrice: total, CreatedAt: now, UpdatedAt: now}
			if p.IsDigital() && item.ProductType == models.TypePhysical {
				item.ProductType = models.TypeDigitalLicense
			}
			order.Items = append(order.Items, item)
			order.TotalAmount += total
			order.HasDigital = order.HasDigital || p.IsDigital()
			order.HasPhysical = order.HasPhysical || !p.IsDigital()
		}
		order.TotalAmount = math.Round(order.TotalAmount*100) / 100
		if err := tx.Create(order); err != nil {
			return err
		}
		for i := range order.Items {
			if err := tx.Create(&order.Items[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
