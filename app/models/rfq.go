package models

import "time"

type RfqStatus string

const (
	RfqStatusOpen    RfqStatus = "open"
	RfqStatusAwarded RfqStatus = "awarded"
	RfqStatusClosed  RfqStatus = "closed"
)

type RfqQuote struct {
	ID            string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	RfqID         string    `gorm:"type:varchar(64);index" json:"rfq_id"`
	VendorID      string    `gorm:"type:varchar(64)" json:"vendor_id"`
	VendorName    string    `gorm:"type:varchar(255)" json:"vendor_name"`
	VendorCountry string    `gorm:"type:varchar(64)" json:"vendor_country"`
	UnitPrice     float64   `gorm:"type:double" json:"unit_price"`
	TotalAmount   float64   `gorm:"type:double" json:"total_amount"`
	LeadTimeDays  int       `gorm:"type:integer" json:"lead_time_days"`
	Notes         string    `gorm:"type:text" json:"notes"`
	Status        string    `gorm:"type:varchar(32)" json:"status"` // pending, accepted, rejected
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (RfqQuote) TableName() string {
	return "rfq_quotes"
}

type Rfq struct {
	ID                 string     `gorm:"primaryKey;type:varchar(64)" json:"id"`
	RfqNumber          string     `gorm:"type:varchar(64);uniqueIndex" json:"rfq_number"`
	BuyerID            string     `gorm:"type:varchar(64);index" json:"buyer_id"`
	BuyerName          string     `gorm:"type:varchar(255)" json:"buyer_name"`
	BuyerCompany       string     `gorm:"type:varchar(255)" json:"buyer_company"`
	BuyerCountry       string     `gorm:"type:varchar(64)" json:"buyer_country"`
	Title              string     `gorm:"type:varchar(255)" json:"title"`
	Category           string     `gorm:"type:varchar(128)" json:"category"`
	Quantity           int        `gorm:"type:integer" json:"quantity"`
	TargetPrice        float64    `gorm:"type:double" json:"target_price"`
	DestinationCountry string     `gorm:"type:varchar(64)" json:"destination_country"`
	Specifications     string     `gorm:"type:text" json:"specifications"`
	Status             RfqStatus  `gorm:"type:varchar(32);index" json:"status"`
	Quotes             []RfqQuote `gorm:"foreignKey:RfqID" json:"quotes"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Rfq) TableName() string {
	return "rfqs"
}
