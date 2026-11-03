package models

import "time"

type OfferStatus string

const (
	OfferStatusPending   OfferStatus = "pending"
	OfferStatusAccepted  OfferStatus = "accepted"
	OfferStatusDeclined  OfferStatus = "declined"
	OfferStatusCountered OfferStatus = "countered"
)

type Offer struct {
	ID           string      `gorm:"primaryKey;type:varchar(64)" json:"id"`
	ProductID    string      `gorm:"type:varchar(64);index" json:"product_id"`
	ProductTitle string      `gorm:"type:varchar(255)" json:"product_title"`
	BuyerID      string      `gorm:"type:varchar(64);index" json:"buyer_id"`
	BuyerName    string      `gorm:"type:varchar(255)" json:"buyer_name"`
	BuyerEmail   string      `gorm:"type:varchar(255)" json:"buyer_email"`
	BuyerCompany string      `gorm:"type:varchar(255)" json:"buyer_company"`
	OfferedPrice float64     `gorm:"type:double" json:"offered_price"`
	Quantity     int         `gorm:"type:integer" json:"quantity"`
	TotalOffered float64     `gorm:"type:double" json:"total_offered"`
	Message      string      `gorm:"type:text" json:"message"`
	Status       OfferStatus `gorm:"type:varchar(32);index" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (Offer) TableName() string {
	return "offers"
}
