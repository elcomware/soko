package dto

import (
	"time"

	"elcomware/app/models"
)

type CheckoutItemInput struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type CheckoutInput struct {
	UserID        string              `json:"user_id"`
	CustomerName  string              `json:"customer_name"`
	CustomerEmail string              `json:"customer_email"`
	Items         []CheckoutItemInput `json:"items"`
	Currency      string              `json:"currency"`
	PaymentMethod string              `json:"payment_method"`
}

type OrderItemOutput struct {
	ID                 string             `json:"id"`
	ProductID          string             `json:"product_id"`
	ProductTitle       string             `json:"product_title"`
	ProductType        models.ProductType `json:"product_type"`
	Quantity           int                `json:"quantity"`
	UnitPrice          float64            `json:"unit_price"`
	TotalPrice         float64            `json:"total_price"`
	DigitalLicenseKey  string             `json:"digital_license_key,omitempty"`
	DigitalDownloadUrl string             `json:"digital_download_url,omitempty"`
}

type OrderOutput struct {
	ID            string            `json:"id"`
	OrderNumber   string            `json:"order_number"`
	UserID        string            `json:"user_id"`
	CustomerName  string            `json:"customer_name"`
	CustomerEmail string            `json:"customer_email"`
	Items         []OrderItemOutput `json:"items"`
	TotalAmount   float64           `json:"total_amount"`
	Currency      string            `json:"currency"`
	Status        string            `json:"status"`
	PaymentMethod string            `json:"payment_method"`
	HasDigital    bool              `json:"has_digital"`
	HasPhysical   bool              `json:"has_physical"`
	CreatedAt     time.Time         `json:"created_at"`
}
