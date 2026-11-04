package queries

import (
	"strings"
	"testing"
)

func TestStorefrontContentValidation(t *testing.T) {
	valid := DefaultStorefrontContent()
	tests := []struct {
		name, key, value string
		wantError        bool
	}{
		{"blank title", "hero_heading", "  ", true},
		{"unicode title at limit", "hero_heading", strings.Repeat("界", 120), false},
		{"unicode title over limit", "hero_heading", strings.Repeat("界", 121), true},
		{"description limit", "hero_description", strings.Repeat("a", 501), true},
		{"announcement optional", "announcement", "", false},
		{"markup rejected", "announcement", "<script>alert(1)</script>", true},
		{"line breaks allowed", "hero_description", "First line\nSecond line", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := DefaultStorefrontContent()
			for key, value := range valid {
				values[key] = value
			}
			values[tt.key] = tt.value
			_, err := ValidateStorefrontContent(values)
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestStorefrontContentMergeOnlyKnownKeys(t *testing.T) {
	content := mergeStorefrontContent(`{"hero_heading":"A better way to source","announcement":"","untrusted_key":"discard"}`)
	if content["hero_heading"] != "A better way to source" {
		t.Fatal("persisted title not applied")
	}
	if content["hero_accent"] != DefaultStorefrontContent()["hero_accent"] {
		t.Fatal("missing fields must retain defaults")
	}
	if content["announcement"] != "" {
		t.Fatal("optional copy must be clearable")
	}
	if _, ok := content["untrusted_key"]; ok {
		t.Fatal("unexpected content key leaked")
	}
	if mergeStorefrontContent(`{invalid`)["hero_heading"] != DefaultStorefrontContent()["hero_heading"] {
		t.Fatal("invalid stored JSON must fall back safely")
	}
	copy := DefaultStorefrontContent()
	copy["hero_heading"] = "changed"
	if DefaultStorefrontContent()["hero_heading"] == "changed" {
		t.Fatal("defaults shared mutable state")
	}
}
