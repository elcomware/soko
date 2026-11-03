package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type CreateRfqRequest struct {
	Title              string  `form:"title" json:"title"`
	Category           string  `form:"category" json:"category"`
	Quantity           int     `form:"quantity" json:"quantity"`
	TargetPrice        float64 `form:"target_price" json:"target_price"`
	DestinationCountry string  `form:"destination_country" json:"destination_country"`
	Specifications     string  `form:"specifications" json:"specifications"`
	BuyerName          string  `form:"buyer_name" json:"buyer_name"`
	BuyerCompany       string  `form:"buyer_company" json:"buyer_company"`
	BuyerCountry       string  `form:"buyer_country" json:"buyer_country"`
}

func (r *CreateRfqRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CreateRfqRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"title":    "required|min_len:5",
		"category": "required",
		"quantity": "required|min:1",
	}
}

func (r *CreateRfqRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"title.required":    "RFQ brief title is required.",
		"title.min_len":     "RFQ brief title must be at least 5 characters.",
		"category.required": "Product category is required.",
		"quantity.required": "Target sourcing quantity is required.",
	}
}

