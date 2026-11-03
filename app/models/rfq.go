package models

import "time"

type RfqStatus string

const (
	RfqStatusOpen    RfqStatus = "open"
	RfqStatusAwarded RfqStatus = "awarded"
	RfqStatusClosed  RfqStatus = "closed"
)

type RfqQuote struct {
	ID           string    `json:"id"`
	RfqID        string    `json:"rfq_id"`
	VendorID     string    `json:"vendor_id"`
	VendorName   string    `json:"vendor_name"`
	VendorCountry string   `json:"vendor_country"`
	UnitPrice    float64   `json:"unit_price"`
	TotalAmount  float64   `json:"total_amount"`
	LeadTimeDays int       `json:"lead_time_days"`
	Notes        string    `json:"notes"`
	Status       string    `json:"status"` // pending, accepted, rejected
	CreatedAt    time.Time `json:"created_at"`
}

type Rfq struct {
	ID                 string     `json:"id"`
	RfqNumber          string     `json:"rfq_number"`
	BuyerID            string     `json:"buyer_id"`
	BuyerName          string     `json:"buyer_name"`
	BuyerCompany       string     `json:"buyer_company"`
	BuyerCountry       string     `json:"buyer_country"`
	Title              string     `json:"title"`
	Category           string     `json:"category"`
	Quantity           int        `json:"quantity"`
	TargetPrice        float64    `json:"target_price"`
	DestinationCountry string     `json:"destination_country"`
	Specifications     string     `json:"specifications"`
	Status             RfqStatus  `json:"status"`
	Quotes             []RfqQuote `json:"quotes"`
	CreatedAt          time.Time  `json:"created_at"`
}

