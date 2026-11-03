package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type SubmitRfqQuoteRequest struct {
	RfqID        string  `form:"rfq_id" json:"rfq_id"`
	UnitPrice    float64 `form:"unit_price" json:"unit_price"`
	LeadTimeDays int     `form:"lead_time_days" json:"lead_time_days"`
	Notes        string  `form:"notes" json:"notes"`
}

func (r *SubmitRfqQuoteRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *SubmitRfqQuoteRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"rfq_id":     "required",
		"unit_price": "required",
	}
}

func (r *SubmitRfqQuoteRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"rfq_id.required":     "RFQ identifier is required.",
		"unit_price.required": "Quotation unit price is required.",
	}
}

