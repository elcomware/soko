package queries

import (
	"encoding/json"
	"testing"
)

func TestCardFooterValidation(t *testing.T) {
	for _, mode := range []string{"together", "rotate"} {
		for _, interval := range []int{5, 8, 60} {
			b := DefaultBrandSettings()
			b.CardFooterMode = mode
			b.CardFooterIntervalSeconds = interval
			b.CardFooterPrice = false
			b.CardFooterRating = false
			b.CardFooterAvailability = false
			if err := ValidateBrandSettings(b); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, interval := range []int{-1, 0, 4, 61} {
		b := DefaultBrandSettings()
		b.CardFooterIntervalSeconds = interval
		if ValidateBrandSettings(b) == nil {
			t.Fatal("accepted invalid interval", interval)
		}
	}
	b := DefaultBrandSettings()
	b.CardFooterMode = "arbitrary"
	if ValidateBrandSettings(b) == nil {
		t.Fatal("accepted invalid mode")
	}
}
func TestLegacyBrandRecordGetsFooterDefaults(t *testing.T) {
	b := DefaultBrandSettings()
	if err := json.Unmarshal([]byte(`{"display_name":"Existing brand","version":7}`), &b); err != nil {
		t.Fatal(err)
	}
	if !b.CardFooterPrice || !b.CardFooterAvailability || b.CardFooterRating || b.CardFooterMode != "together" || b.CardFooterIntervalSeconds != 8 {
		t.Fatalf("unexpected defaults: %+v", b)
	}
	if err := json.Unmarshal([]byte(`{"card_footer_price":false,"card_footer_availability":false}`), &b); err != nil {
		t.Fatal(err)
	}
	if b.CardFooterPrice || b.CardFooterAvailability {
		t.Fatal("explicit false must override defaults")
	}
}
