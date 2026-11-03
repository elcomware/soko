package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusShipped   OrderStatus = "shipped"
)

type OrderItem struct {
	ProductID          string      `json:"product_id"`
	ProductTitle       string      `json:"product_title"`
	ProductType        ProductType `json:"product_type"`
	Quantity           int         `json:"quantity"`
	UnitPrice          float64     `json:"unit_price"`
	TotalPrice         float64     `json:"total_price"`
	DigitalLicenseKey  string      `json:"digital_license_key,omitempty"`
	DigitalDownloadUrl string      `json:"digital_download_url,omitempty"`
}

type Order struct {
	ID            string      `json:"id"`
	OrderNumber   string      `json:"order_number"`
	UserID        string      `json:"user_id"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	Items         []OrderItem `json:"items"`
	TotalAmount   float64     `json:"total_amount"`
	Currency      string      `json:"currency"`
	Status        OrderStatus `json:"status"`
	PaymentMethod string      `json:"payment_method"`
	HasDigital    bool        `json:"has_digital"`
	HasPhysical   bool        `json:"has_physical"`
	CreatedAt     time.Time   `json:"created_at"`
}

