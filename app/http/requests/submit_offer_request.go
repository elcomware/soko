package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type SubmitOfferRequest struct {
	ProductID    string  `form:"product_id" json:"product_id"`
	OfferedPrice float64 `form:"offered_price" json:"offered_price"`
	Quantity     int     `form:"quantity" json:"quantity"`
	Name         string  `form:"name" json:"name"`
	Email        string  `form:"email" json:"email"`
	Company      string  `form:"company" json:"company"`
	Message      string  `form:"message" json:"message"`
}

func (r *SubmitOfferRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *SubmitOfferRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"product_id":    "required",
		"offered_price": "required",
		"quantity":      "required|min:1",
	}
}

func (r *SubmitOfferRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"product_id.required":    "Target product ID is required.",
		"offered_price.required": "Proposed offer price is required.",
		"quantity.required":      "Offer quantity is required.",
	}
}
