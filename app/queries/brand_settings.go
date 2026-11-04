package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type BrandSettings struct {
	DisplayName               string `json:"display_name"`
	Tagline                   string `json:"tagline"`
	PrimaryColor              string `json:"primary_color"`
	AccentColor               string `json:"accent_color"`
	FontBody                  string `json:"font_body"`
	FontDisplay               string `json:"font_display"`
	Version                   int    `json:"version"`
	CardFooterPrice           bool   `json:"card_footer_price"`
	CardFooterRating          bool   `json:"card_footer_rating"`
	CardFooterAvailability    bool   `json:"card_footer_availability"`
	CardFooterMode            string `json:"card_footer_mode"`
	CardFooterIntervalSeconds int    `json:"card_footer_interval_seconds"`
}

func DefaultBrandSettings() BrandSettings {
	return BrandSettings{DisplayName: "Soko", Tagline: "Global commerce", PrimaryColor: "#142f48", AccentColor: "#f3ac32", FontBody: "Manrope", FontDisplay: "Manrope", CardFooterPrice: true, CardFooterAvailability: true, CardFooterMode: "together", CardFooterIntervalSeconds: 8}
}
func LoadBrandSettings() BrandSettings {
	b := DefaultBrandSettings()
	if !facades.Schema().HasTable("brand_settings") {
		return b
	}
	var r models.BrandSettingsRecord
	if facades.Orm().Query().Where("id = ?", "brand").First(&r) != nil || r.ID == "" {
		return b
	}
	if json.Unmarshal([]byte(r.SettingsJSON), &b) != nil {
		return DefaultBrandSettings()
	}
	b.Version = r.Version
	return b
}
func ValidateBrandSettings(b BrandSettings) error {
	if b.CardFooterMode != "together" && b.CardFooterMode != "rotate" {
		return errors.New("choose together or rotate for card details")
	}
	if b.CardFooterIntervalSeconds < 5 || b.CardFooterIntervalSeconds > 60 {
		return errors.New("card rotation interval must be between 5 and 60 seconds")
	}
	if len(strings.TrimSpace(b.DisplayName)) < 1 || len(b.DisplayName) > 50 || len(b.Tagline) > 100 {
		return errors.New("name or tagline is too long or empty")
	}
	hex := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	if !hex.MatchString(b.PrimaryColor) || !hex.MatchString(b.AccentColor) {
		return errors.New("use six-digit hex colors")
	}
	for _, f := range []string{b.FontBody, b.FontDisplay} {
		if f != "Manrope" && f != "system" {
			return errors.New("choose a bundled font")
		}
	}
	return nil
}
