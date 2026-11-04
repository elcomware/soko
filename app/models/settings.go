package models

type BrandSettingsRecord struct {
	ID           string `gorm:"primaryKey;type:varchar(32)"`
	SettingsJSON string
	Version      int
}

func (BrandSettingsRecord) TableName() string { return "brand_settings" }

type MerchantPaymentSettings struct {
	MerchantID   string `gorm:"primaryKey;type:varchar(64)"`
	SettingsJSON string
	Version      int
}

func (MerchantPaymentSettings) TableName() string { return "merchant_payment_settings" }
