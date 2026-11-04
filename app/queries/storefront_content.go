package queries

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"elcomware/app/facades"
	"elcomware/app/models"
)

// DefaultStorefrontContent returns a fresh map, safe for callers to modify.
func DefaultStorefrontContent() map[string]string {
	return map[string]string{
		"hero_heading":         "The world’s supply.",
		"hero_accent":          "On your terms.",
		"hero_description":     "From factory floors to your next frontier. Discover equipment, connect with suppliers and build a business without boundaries.",
		"announcement":         "A world of suppliers. One marketplace.",
		"sourcing_heading":     "Your next big thing starts with a brief.",
		"sourcing_description": "Some needs don’t fit a search box. Share your specifications and give suppliers a clear place to start.",
	}
}

func ValidateStorefrontContent(values map[string]string) (map[string]string, error) {
	result := make(map[string]string, 6)
	fields := []struct {
		key, label string
		max        int
		required   bool
	}{
		{"hero_heading", "Hero heading", 120, true},
		{"hero_accent", "Hero accent", 120, true},
		{"hero_description", "Hero description", 500, false},
		{"announcement", "Announcement", 160, false},
		{"sourcing_heading", "Sourcing heading", 120, true},
		{"sourcing_description", "Sourcing description", 500, false},
	}
	for _, field := range fields {
		value := strings.TrimSpace(values[field.key])
		if !utf8.ValidString(value) || utf8.RuneCountInString(value) > field.max {
			return nil, fmt.Errorf("%s must be at most %d characters", field.label, field.max)
		}
		if field.required && value == "" {
			return nil, fmt.Errorf("%s is required", field.label)
		}
		if strings.ContainsAny(value, "<>") {
			return nil, fmt.Errorf("%s must contain plain text only", field.label)
		}
		for _, char := range value {
			if unicode.IsControl(char) && char != '\n' && char != '\t' {
				return nil, fmt.Errorf("%s contains unsupported characters", field.label)
			}
		}
		result[field.key] = value
	}
	return result, nil
}

func mergeStorefrontContent(raw string) map[string]string {
	result := DefaultStorefrontContent()
	var saved map[string]string
	if json.Unmarshal([]byte(raw), &saved) != nil {
		return result
	}
	// Validate persisted values too; corrupted or legacy data cannot break the page.
	combined := DefaultStorefrontContent()
	for key := range combined {
		if value, ok := saved[key]; ok {
			combined[key] = value
		}
	}
	valid, err := ValidateStorefrontContent(combined)
	if err != nil {
		return result
	}
	return valid
}

// LoadStorefrontContent remains usable before the additive migration is applied.
func LoadStorefrontContent() map[string]string {
	if !facades.Schema().HasTable("storefront_contents") {
		return DefaultStorefrontContent()
	}
	var content models.StorefrontContent
	if err := facades.Orm().Query().Where("id = ?", "store").First(&content); err != nil || content.ID == "" {
		return DefaultStorefrontContent()
	}
	return mergeStorefrontContent(content.ContentJSON)
}
