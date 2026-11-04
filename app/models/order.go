package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusShipped   OrderStatus = "shipped"
)

type OrderItem struct {
	OriginalUnitPrice  float64     `json:"original_unit_price"`
	DiscountPerUnit    float64     `json:"discount_per_unit"`
	PromotionID        string      `gorm:"type:varchar(64)" json:"promotion_id,omitempty"`
	PromotionName      string      `gorm:"type:varchar(120)" json:"promotion_name,omitempty"`
	StockReserved      bool        `gorm:"type:boolean;default:false" json:"-"`
	ID                 string      `gorm:"primaryKey;type:varchar(64)" json:"id"`
	OrderID            string      `gorm:"type:varchar(64);index" json:"order_id"`
	ProductID          string      `gorm:"type:varchar(64);index" json:"product_id"`
	ProductTitle       string      `gorm:"type:varchar(255)" json:"product_title"`
	ProductType        ProductType `gorm:"type:varchar(64)" json:"product_type"`
	Quantity           int         `gorm:"type:integer" json:"quantity"`
	UnitPrice          float64     `gorm:"type:double" json:"unit_price"`
	TotalPrice         float64     `gorm:"type:double" json:"total_price"`
	DigitalLicenseKey  string      `gorm:"type:text" json:"digital_license_key,omitempty"`
	DigitalDownloadUrl string      `gorm:"type:text" json:"digital_download_url,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type Order struct {
	ID            string      `gorm:"primaryKey;type:varchar(64)" json:"id"`
	OrderNumber   string      `gorm:"type:varchar(64);uniqueIndex" json:"order_number"`
	UserID        string      `gorm:"type:varchar(64);index" json:"user_id"`
	CustomerName  string      `gorm:"type:varchar(255)" json:"customer_name"`
	CustomerEmail string      `gorm:"type:varchar(255)" json:"customer_email"`
	Items         []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	TotalAmount   float64     `gorm:"type:double" json:"total_amount"`
	Currency      string      `gorm:"type:varchar(16)" json:"currency"`
	Status        OrderStatus `gorm:"type:varchar(32)" json:"status"`
	PaymentMethod string      `gorm:"type:varchar(64)" json:"payment_method"`
	HasDigital    bool        `gorm:"type:boolean" json:"has_digital"`
	HasPhysical   bool        `gorm:"type:boolean" json:"has_physical"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (Order) TableName() string {
	return "orders"
}
