package dto

import "time"

type CreateRfqInput struct {
	BuyerID            string  `json:"buyer_id"`
	BuyerName          string  `json:"buyer_name"`
	BuyerCompany       string  `json:"buyer_company"`
	BuyerCountry       string  `json:"buyer_country"`
	Title              string  `json:"title"`
	Category           string  `json:"category"`
	Quantity           int     `json:"quantity"`
	TargetPrice        float64 `json:"target_price"`
	DestinationCountry string  `json:"destination_country"`
	Specifications     string  `json:"specifications"`
}

type SubmitQuoteInput struct {
	RfqID        string  `json:"rfq_id"`
	VendorID     string  `json:"vendor_id"`
	VendorName   string  `json:"vendor_name"`
	VendorCountry string `json:"vendor_country"`
	UnitPrice    float64 `json:"unit_price"`
	LeadTimeDays int     `json:"lead_time_days"`
	Notes        string  `json:"notes"`
}

type RfqQuoteOutput struct {
	ID            string    `json:"id"`
	RfqID         string    `json:"rfq_id"`
	VendorID      string    `json:"vendor_id"`
	VendorName    string    `json:"vendor_name"`
	VendorCountry string    `json:"vendor_country"`
	UnitPrice     float64   `json:"unit_price"`
	TotalAmount   float64   `json:"total_amount"`
	LeadTimeDays  int       `json:"lead_time_days"`
	Notes         string    `json:"notes"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type RfqOutput struct {
	ID                 string           `json:"id"`
	RfqNumber          string           `json:"rfq_number"`
	BuyerID            string           `json:"buyer_id"`
	BuyerName          string           `json:"buyer_name"`
	BuyerCompany       string           `json:"buyer_company"`
	BuyerCountry       string           `json:"buyer_country"`
	Title              string           `json:"title"`
	Category           string           `json:"category"`
	Quantity           int              `json:"quantity"`
	TargetPrice        float64          `json:"target_price"`
	DestinationCountry string           `json:"destination_country"`
	Specifications     string           `json:"specifications"`
	Status             string           `json:"status"`
	Quotes             []RfqQuoteOutput `json:"quotes"`
	CreatedAt          time.Time        `json:"created_at"`
}
