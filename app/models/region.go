package models

import (
	"strings"
	"unicode"
)

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

// NormalizeRegion accepts current IDs and the labels used by existing catalogs.
// It does not rewrite persisted rows when producing a public or workspace DTO.
func NormalizeRegion(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return "global"
	}
	value = strings.Join(strings.FieldsFunc(value, func(r rune) bool { return r == '_' || r == '-' || r == '–' || r == '—' || unicode.IsSpace(r) }), "-")
	switch value {
	case "global-unified-marketplace", "universal-unified-marketplace":
		return "global"
	case "east-africa-(eac-hub)", "east-africa-(eac-corridor)":
		return "east-africa"
	case "west-africa-(ecowas-hub)":
		return "west-africa"
	case "north-america-(us-&-canada)":
		return "north-america"
	case "europe-&-uk", "european-union-&-uk":
		return "europe"
	case "asia-pacific-(shenzhen-&-hk-hub)":
		return "asia-pacific"
	case "middle-east-(gcc-hub)":
		return "middle-east"
	}
	return value
}

// RegionAliases includes the canonical, former underscore and friendly DB values.
func RegionAliases(raw string) []string {
	canonical := NormalizeRegion(raw)
	aliases := []string{canonical, strings.ReplaceAll(canonical, "-", "_"), strings.ReplaceAll(canonical, "-", " ")}
	extras := map[string][]string{
		"global":        {"global unified marketplace", "universal unified marketplace"},
		"east-africa":   {"east africa (eac hub)", "east africa (eac corridor)"},
		"west-africa":   {"west africa (ecowas hub)"},
		"north-america": {"north america (us & canada)"},
		"europe":        {"europe & uk", "european union & uk"},
		"asia-pacific":  {"asia-pacific (shenzhen & hk hub)"},
		"middle-east":   {"middle east (gcc hub)"},
	}
	return append(aliases, extras[canonical]...)
}
