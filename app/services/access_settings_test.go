package services

import (
	"elcomware/app/models"
	"elcomware/app/queries"
	"testing"
)

func TestValidateGrants(t *testing.T) {
	for _, g := range [][]models.PermissionGrant{{{"unknown", "all"}}, {{"products.read", "everything"}}, {{"permissions.manage", "own"}}, {{"products.read", "own"}, {"products.read", "all"}}} {
		if ValidateGrants(g) == nil {
			t.Fatal("accepted invalid grants")
		}
	}
	if ValidateGrants([]models.PermissionGrant{{"products.read", "own"}}) != nil {
		t.Fatal("valid rejected")
	}
}
func TestPaymentPreferences(t *testing.T) {
	p := PaymentPreference{Provider: "stripe", Methods: []string{"card"}, Country: "US", Currency: "USD"}
	if ValidatePaymentPreference(p) != nil {
		t.Fatal("valid preferences rejected")
	}
	p.Enabled = true
	if ValidatePaymentPreference(p) == nil {
		t.Fatal("live collection enabled without provider")
	}
	p.Enabled = false
	p.Methods = []string{"mobile_money"}
	if ValidatePaymentPreference(p) == nil {
		t.Fatal("unsupported method accepted")
	}
}
func TestBrandValidation(t *testing.T) {
	b := queries.DefaultBrandSettings()
	if queries.ValidateBrandSettings(b) != nil {
		t.Fatal("defaults invalid")
	}
	b.PrimaryColor = "url(javascript:alert(1))"
	if queries.ValidateBrandSettings(b) == nil {
		t.Fatal("unsafe color")
	}
}
