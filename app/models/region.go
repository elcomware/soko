package models

type Currency struct {
	Code         string  `json:"code"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	ExchangeRate float64 `json:"exchange_rate"` // Relative to USD = 1.0
	Prefix       bool    `json:"prefix"`        // True for $100, False for 100 د.إ
}

type MarketRegion struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Flag               string   `json:"flag"`
	DefaultCurrency    string   `json:"default_currency"`
	LocalLogistics     string   `json:"local_logistics"`
	PaymentRails       []string `json:"payment_rails"`
	VatPercentage      float64  `json:"vat_percentage"`
	CustomsClearanceHr int      `json:"customs_clearance_hr"`
}

