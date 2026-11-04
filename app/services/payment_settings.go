package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/payments"
	"encoding/json"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"regexp"
)

// PaymentPreference stores merchant choices, never credentials or a claim of provider readiness.
type PaymentPreference struct {
	Provider string                   `json:"provider"`
	Methods  []string                 `json:"methods"`
	Currency string                   `json:"currency"`
	Country  string                   `json:"country"`
	Enabled  bool                     `json:"enabled"`
	Version  int                      `json:"version"`
	Status   string                   `json:"status"`
	Routes   []PaymentRoutePreference `json:"routes"`
}

func DefaultPaymentPreference() PaymentPreference {
	return PaymentPreference{Provider: "none", Methods: []string{}, Currency: "USD", Status: "unconfigured", Routes: []PaymentRoutePreference{}}
}
func LoadPaymentPreference(merchant string) PaymentPreference {
	p := DefaultPaymentPreference()
	if !facades.Schema().HasTable("merchant_payment_settings") {
		return p
	}
	var r models.MerchantPaymentSettings
	if facades.Orm().Query().Where("merchant_id = ?", merchant).First(&r) != nil || r.MerchantID == "" {
		return p
	}
	if json.Unmarshal([]byte(r.SettingsJSON), &p) != nil {
		return DefaultPaymentPreference()
	}
	p.Version = r.Version
	p.Enabled = false
	p.Status = "requires_provider_onboarding"
	if p.Provider == "none" && len(p.Routes) == 0 {
		p.Status = "unconfigured"
	}
	return p
}
func ValidatePaymentPreference(p PaymentPreference) error {
	if err := ValidatePaymentRoutes(p.Routes); err != nil {
		return err
	}
	allowed := map[string]bool{"none": true, "stripe": true, "flutterwave": true, "manual_bank_transfer": true}
	if !allowed[p.Provider] {
		return errors.New("unknown payment provider")
	}
	if p.Enabled {
		return errors.New("live collection cannot be enabled until server-side provider onboarding and credentials are verified")
	}
	if len(p.Country) != 2 || len(p.Currency) != 3 {
		return errors.New("use two-letter country and three-letter currency codes")
	}
	for _, s := range []string{p.Country, p.Currency} {
		for _, r := range s {
			if r < 'A' || r > 'Z' {
				return errors.New("country and currency codes must be uppercase letters")
			}
		}
	}
	seen := map[string]bool{}
	for _, m := range p.Methods {
		if seen[m] || !map[string]bool{"card": true, "mobile_money": true, "bank_transfer": true}[m] {
			return errors.New("invalid payment method")
		}
		seen[m] = true
	}
	if p.Provider == "manual_bank_transfer" && (len(p.Methods) != 1 || p.Methods[0] != "bank_transfer") {
		return errors.New("manual transfer requires bank transfer only")
	}
	if p.Provider == "stripe" && seen["mobile_money"] {
		return errors.New("mobile money is not enabled for this Stripe adapter")
	}
	return nil
}
func SavePaymentPreference(actor *models.User, merchant string, p PaymentPreference) error {
	if !actor.Can("payments.configure", merchant) {
		return errors.New("payment configuration access denied")
	}
	if err := ValidatePaymentPreference(p); err != nil {
		return err
	}
	p.Status = "requires_provider_onboarding"
	if p.Provider == "none" && len(p.Routes) == 0 {
		p.Status = "unconfigured"
	}
	data, _ := json.Marshal(p)
	return facades.Orm().Transaction(func(tx orm.Query) error {
		var r models.MerchantPaymentSettings
		if err := tx.Where("merchant_id = ?", merchant).First(&r); err != nil {
			return err
		}
		if r.Version != p.Version {
			return errors.New("payment settings changed; refresh")
		}
		if r.MerchantID == "" {
			return tx.Create(&models.MerchantPaymentSettings{MerchantID: merchant, SettingsJSON: string(data), Version: 1})
		}
		res, err := tx.Model(&models.MerchantPaymentSettings{}).Where("merchant_id = ? AND version = ?", merchant, p.Version).Update(map[string]any{"settings_json": string(data), "version": p.Version + 1})
		if err == nil && res.RowsAffected != 1 {
			return errors.New("payment settings changed")
		}
		return err
	})
}

// PaymentRoutePreference is an unverified requested route, never a capability grant.
// No credentials, account identifiers or live enable switch belong in this record.
type PaymentRoutePreference struct {
	Country   string   `json:"country"`
	Area      string   `json:"area"`
	Currency  string   `json:"currency"`
	Method    string   `json:"method"`
	Providers []string `json:"providers"`
}

func ValidatePaymentRoutes(routes []PaymentRoutePreference) error {
	if len(routes) > 50 {
		return errors.New("at most 50 payment routes are supported")
	}
	seen := map[string]bool{}
	code := regexp.MustCompile(`^[A-Z]{2}$`)
	currency := regexp.MustCompile(`^[A-Z]{3}$`)
	area := regexp.MustCompile(`^[A-Z0-9-]{0,32}$`)
	for _, r := range routes {
		if !code.MatchString(r.Country) || !currency.MatchString(r.Currency) || !area.MatchString(r.Area) {
			return errors.New("routes require uppercase country/currency codes and an optional area code of up to 32 letters, digits or hyphens")
		}
		key := r.Country + ":" + r.Area + ":" + r.Currency + ":" + r.Method
		if seen[key] {
			return errors.New("duplicate payment market route")
		}
		seen[key] = true
		if len(r.Providers) < 1 || len(r.Providers) > 3 {
			return errors.New("choose one to three providers per route")
		}
		providers := map[string]bool{}
		for _, id := range r.Providers {
			if providers[id] {
				return errors.New("duplicate provider in payment route")
			}
			providers[id] = true
			switch id {
			case "stripe":
				if r.Method != "card" && r.Method != "bank_transfer" {
					return errors.New("requested method is not supported by the Stripe contract")
				}
			case "flutterwave":
				if r.Method != "card" && r.Method != "bank_transfer" && r.Method != "mobile_money" {
					return errors.New("unknown requested payment method")
				}
			case "manual_bank_transfer":
				if r.Method != "bank_transfer" {
					return errors.New("manual transfer routes require bank_transfer")
				}
			default:
				return errors.New("unknown route provider")
			}
		}
	}
	return nil
}

// SelectPaymentRoute prefers an exact area rule over a country-wide rule.
// It only resolves merchant intent; callers MUST run payments.Candidates using
// server-verified adapters before execution. No route activates collection.
func SelectPaymentRoute(merchant string, market payments.Market, routes []PaymentRoutePreference) (payments.Route, error) {
	if merchant == "" {
		return payments.Route{}, payments.ErrNoRoute
	}
	if err := ValidatePaymentRoutes(routes); err != nil {
		return payments.Route{}, err
	}
	var selected *PaymentRoutePreference
	for i := range routes {
		r := &routes[i]
		if r.Country != market.Country || r.Currency != market.Currency || r.Method != market.Method {
			continue
		}
		if r.Area == market.Area {
			selected = r
			break
		}
		if r.Area == "" {
			selected = r
		}
	}
	if selected == nil {
		return payments.Route{}, payments.ErrNoRoute
	}
	return payments.Route{MerchantID: merchant, Market: payments.Market{Country: selected.Country, Area: selected.Area, Currency: selected.Currency, Method: selected.Method}, Providers: append([]string(nil), selected.Providers...)}, nil
}
