package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type CheckoutRequest struct {
	ProductID     string  `form:"product_id" json:"product_id"`
	Quantity      int     `form:"quantity" json:"quantity"`
	CustomerName  string  `form:"customer_name" json:"customer_name"`
	CustomerEmail string  `form:"customer_email" json:"customer_email"`
	Currency      string  `form:"currency" json:"currency"`
	PaymentMethod string  `form:"payment_method" json:"payment_method"`
	UnitPrice     float64 `form:"unit_price" json:"unit_price"`
}

func (r *CheckoutRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CheckoutRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"product_id": "required",
		"quantity":   "required|min:1",
	}
}

func (r *CheckoutRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"product_id.required": "Product selection is required.",
		"quantity.required":   "Order quantity is required.",
		"quantity.min":        "Order quantity must be at least 1 unit.",
	}
}
