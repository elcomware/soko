package dto

import "time"

type SubmitOfferInput struct {
	ProductID    string  `json:"product_id"`
	BuyerID      string  `json:"buyer_id"`
	BuyerName    string  `json:"buyer_name"`
	BuyerEmail   string  `json:"buyer_email"`
	BuyerCompany string  `json:"buyer_company"`
	OfferedPrice float64 `json:"offered_price"`
	Quantity     int     `json:"quantity"`
	Message      string  `json:"message"`
}

type OfferOutput struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	ProductTitle string    `json:"product_title"`
	BuyerID      string    `json:"buyer_id"`
	BuyerName    string    `json:"buyer_name"`
	BuyerEmail   string    `json:"buyer_email"`
	BuyerCompany string    `json:"buyer_company"`
	OfferedPrice float64   `json:"offered_price"`
	Quantity     int       `json:"quantity"`
	TotalOffered float64   `json:"total_offered"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
