package models

// TestPaymentAttempt is deliberately separate from orders and fulfillment.
// It stores no credentials, card data or raw webhook bodies.
type TestPaymentAttempt struct {
	ID               string `gorm:"primaryKey;type:varchar(64)"`
	Provider         string
	MerchantID       string
	OrderID          string
	PaymentReference string
	Currency         string
	AmountMinor      int64
	Status           string
	Version          int
}

func (TestPaymentAttempt) TableName() string { return "test_payment_attempts" }

type TestPaymentReceipt struct {
	ID        string `gorm:"primaryKey;type:varchar(64)"` // SHA256(provider + event ID), unique across attempts.
	AttemptID string
	EventHash string
}

func (TestPaymentReceipt) TableName() string { return "test_payment_receipts" }
