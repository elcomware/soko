package actions

import (
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
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
	return a.ExecuteCartWithToken(cart, user, "")
}

// ExecuteCartWithToken makes retries resolve to one unique order and reservation.
func (a *ExecuteCheckoutAction) ExecuteCartWithToken(cart map[string]int, user *models.User, token string) (*models.Order, error) {
	if user == nil || user.ID == "" {
		return nil, errors.New("sign in before placing an order")
	}
	id := uuid.NewString()
	if token != "" {
		if _, err := uuid.Parse(token); err != nil {
			return nil, errors.New("invalid checkout token")
		}
		id = uuid.NewSHA1(uuid.NameSpaceOID, []byte(user.ID+":"+token)).String()
		var existing models.Order
		if err := facades.Orm().Query().Where("id = ? AND user_id = ?", "ord_"+id, user.ID).First(&existing); err != nil {
			return nil, err
		}
		if existing.ID != "" {
			return &existing, nil
		}
	}
	if len(cart) == 0 || len(cart) > 100 {
		return nil, errors.New("cart must contain between 1 and 100 listings")
	}
	now := time.Now()
	order := &models.Order{ID: "ord_" + id, OrderNumber: "SOKO-" + id, UserID: user.ID, CustomerName: user.Name, CustomerEmail: user.Email, Currency: "USD", Status: models.OrderStatusPending, PaymentMethod: "Payment not collected", CreatedAt: now, UpdatedAt: now, Items: []models.OrderItem{}}
	ids := make([]string, 0, len(cart))
	for id := range cart {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		campaigns, err := queries.ActivePromotions(tx, now)
		if err != nil {
			return err
		}
		for _, id := range ids {
			var p models.Product
			if err := tx.Where("id = ?", id).First(&p); err != nil {
				return err
			}
			if p.ID == "" {
				return errors.New("a cart listing is no longer available")
			}
			qty := cart[id]
			pricing, err := models.PriceProduct(&p, qty, campaigns, now)
			price := pricing.UnitPrice
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
			item := models.OrderItem{OriginalUnitPrice: pricing.OriginalUnitPrice, DiscountPerUnit: pricing.DiscountPerUnit, PromotionID: pricing.PromotionID, PromotionName: pricing.PromotionName, StockReserved: p.Stock >= 0, ID: "item_" + uuid.NewString(), OrderID: order.ID, ProductID: id, ProductTitle: p.Title, ProductType: p.Type, Quantity: qty, UnitPrice: price, TotalPrice: total, CreatedAt: now, UpdatedAt: now}
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
		// A concurrent request with the same token may have committed first.
		// Our transaction rolled back its stock changes, so return its order.
		if token != "" {
			var existing models.Order
			if lookupErr := facades.Orm().Query().Where("id = ? AND user_id = ?", order.ID, user.ID).First(&existing); lookupErr == nil && existing.ID != "" {
				return &existing, nil
			}
		}
		return nil, err
	}
	return order, nil
}
