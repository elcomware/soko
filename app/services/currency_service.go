package services

import (
	"elcomware/app/models"
)

type CurrencyService struct {
	currencies map[string]models.Currency
	regions    map[string]models.MarketRegion
}

func NewCurrencyService() *CurrencyService {
	s := &CurrencyService{
		currencies: make(map[string]models.Currency),
		regions:    make(map[string]models.MarketRegion),
	}
	s.initialize()
	return s
}

func (s *CurrencyService) initialize() {
	currencies := []models.Currency{
		{Code: "USD", Symbol: "$", Name: "US Dollar", ExchangeRate: 1.0, Prefix: true},
		{Code: "EUR", Symbol: "€", Name: "Euro", ExchangeRate: 0.92, Prefix: true},
		{Code: "GBP", Symbol: "£", Name: "British Pound", ExchangeRate: 0.78, Prefix: true},
		{Code: "KES", Symbol: "KSh", Name: "Kenya Shilling", ExchangeRate: 130.0, Prefix: false},
		{Code: "NGN", Symbol: "₦", Name: "Nigerian Naira", ExchangeRate: 1550.0, Prefix: true},
		{Code: "ZAR", Symbol: "R", Name: "South African Rand", ExchangeRate: 18.5, Prefix: true},
		{Code: "AED", Symbol: "AED", Name: "UAE Dirham", ExchangeRate: 3.67, Prefix: false},
		{Code: "CNY", Symbol: "¥", Name: "Chinese Yuan", ExchangeRate: 7.23, Prefix: true},
		{Code: "JPY", Symbol: "¥", Name: "Japanese Yen", ExchangeRate: 152.0, Prefix: true},
		{Code: "INR", Symbol: "₹", Name: "Indian Rupee", ExchangeRate: 84.0, Prefix: true},
	}
	for _, c := range currencies {
		s.currencies[c.Code] = c
	}

	regions := []models.MarketRegion{
		{
			ID:                 "global",
			Name:               "Global Unified Marketplace",
			Flag:               "🌐",
			DefaultCurrency:    "USD",
			LocalLogistics:     "Universal multi-currency clearing and inter-continental logistics",
			PaymentRails:       []string{"Trade Assurance Escrow", "Global Swift Wire", "USDC / USDT"},
			VatPercentage:      0.0,
			CustomsClearanceHr: 24,
		},
		{
			ID:                 "east-africa",
			Name:               "East Africa (EAC Hub)",
			Flag:               "🇰🇪",
			DefaultCurrency:    "KES",
			LocalLogistics:     "Kenya, Tanzania, Uganda, Rwanda with M-Pesa & EAC Single Customs",
			PaymentRails:       []string{"M-Pesa B2B Paybill", "EAC Regional RTGS", "Bank Wire"},
			VatPercentage:      16.0,
			CustomsClearanceHr: 12,
		},
		{
			ID:                 "west-africa",
			Name:               "West Africa (ECOWAS Hub)",
			Flag:               "🇳🇬",
			DefaultCurrency:    "NGN",
			LocalLogistics:     "Nigeria, Ghana, Côte d'Ivoire clearing with local escrow settlement",
			PaymentRails:       []string{"NIBSS Instant Payment", "Flutterwave Escrow", "Commercial Wire"},
			VatPercentage:      7.5,
			CustomsClearanceHr: 48,
		},
		{
			ID:                 "north-america",
			Name:               "North America (US & Canada)",
			Flag:               "🇺🇸",
			DefaultCurrency:    "USD",
			LocalLogistics:     "Domestic 2-day dispatch, UL/FCC verified suppliers",
			PaymentRails:       []string{"Fedwire", "ACH Business Transfer", "Trade Credit Line"},
			VatPercentage:      8.25,
			CustomsClearanceHr: 6,
		},
		{
			ID:                 "europe",
			Name:               "European Union & UK",
			Flag:               "🇪🇺",
			DefaultCurrency:    "EUR",
			LocalLogistics:     "CE compliant, VAT automated invoicing, EUR/GBP settlement",
			PaymentRails:       []string{"SEPA Instant Credit Transfer", "BACS UK", "Corporate Credit Card"},
			VatPercentage:      20.0,
			CustomsClearanceHr: 8,
		},
		{
			ID:                 "asia-pacific",
			Name:               "Asia-Pacific (Shenzhen & HK Hub)",
			Flag:               "🇨🇳",
			DefaultCurrency:    "CNY",
			LocalLogistics:     "Direct Tier-1 OEM manufacturing, FOB/CIF bulk container dispatch",
			PaymentRails:       []string{"CIPS Cross-Border RMB", "Alipay+ Merchant", "Telegraphic Transfer"},
			VatPercentage:      13.0,
			CustomsClearanceHr: 18,
		},
		{
			ID:                 "middle-east",
			Name:               "Middle East (GCC Hub)",
			Flag:               "🇦🇪",
			DefaultCurrency:    "AED",
			LocalLogistics:     "Dubai free zone transshipment & GCC regional distribution",
			PaymentRails:       []string{"UAE Central Bank Aani", "SARIE", "Trade Escrow Guarantee"},
			VatPercentage:      5.0,
			CustomsClearanceHr: 4,
		},
	}
	for _, r := range regions {
		s.regions[r.ID] = r
	}
}

func (s *CurrencyService) GetCurrencies() []models.Currency {
	list := make([]models.Currency, 0, len(s.currencies))
	order := []string{"USD", "EUR", "GBP", "KES", "NGN", "ZAR", "AED", "CNY", "JPY", "INR"}
	for _, code := range order {
		if c, ok := s.currencies[code]; ok {
			list = append(list, c)
		}
	}
	return list
}

func (s *CurrencyService) GetRegions() []models.MarketRegion {
	order := []string{"global", "east-africa", "west-africa", "north-america", "europe", "asia-pacific", "middle-east"}
	list := make([]models.MarketRegion, 0, len(order))
	for _, code := range order {
		if r, ok := s.regions[code]; ok {
			list = append(list, r)
		}
	}
	return list
}
