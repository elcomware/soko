package services

import (
	"context"
	"elcomware/app/payments"
	"testing"
)

func TestPaymentRouteConstraints(t *testing.T) {
	valid := PaymentRoutePreference{Country: "KE", Area: "KE-30", Currency: "KES", Method: "mobile_money", Providers: []string{"flutterwave"}}
	if err := ValidatePaymentRoutes([]PaymentRoutePreference{valid}); err != nil {
		t.Fatal(err)
	}
	invalid := []PaymentRoutePreference{
		{Country: "KE", Currency: "KES", Method: "mobile_money", Providers: []string{"stripe"}},
		{Country: "KE", Currency: "KES", Method: "card", Providers: []string{"manual_bank_transfer"}},
		{Country: "KE", Currency: "KES", Method: "card", Providers: []string{"unknown"}},
		{Country: "KE", Currency: "KES", Method: "card", Providers: []string{"stripe", "stripe"}},
		{Country: "KE", Currency: "KES", Method: "card"},
		{Country: "KE", Currency: "KES", Area: "<script>", Method: "card", Providers: []string{"stripe"}},
	}
	for _, r := range invalid {
		if ValidatePaymentRoutes([]PaymentRoutePreference{r}) == nil {
			t.Fatalf("accepted invalid route %+v", r)
		}
	}
	if ValidatePaymentRoutes([]PaymentRoutePreference{valid, valid}) == nil {
		t.Fatal("duplicate route accepted")
	}
	routes := make([]PaymentRoutePreference, 51)
	if ValidatePaymentRoutes(routes) == nil {
		t.Fatal("too many routes accepted")
	}
}
func TestPaymentRouteSelectionDoesNotActivateCollection(t *testing.T) {
	routes := []PaymentRoutePreference{
		{Country: "KE", Currency: "KES", Method: "card", Providers: []string{"stripe", "flutterwave"}},
		{Country: "KE", Area: "KE-30", Currency: "KES", Method: "card", Providers: []string{"flutterwave"}},
	}
	market := payments.Market{Country: "KE", Area: "KE-30", Currency: "KES", Method: "card"}
	r, err := SelectPaymentRoute("merchant", market, routes)
	if err != nil || r.Providers[0] != "flutterwave" || r.MerchantID != "merchant" {
		t.Fatal(r, err)
	}
	if _, err := payments.Candidates(context.Background(), "merchant", market, r, nil); err == nil {
		t.Fatal("draft enabled collection")
	}
	market.Area = "KE-01"
	r, err = SelectPaymentRoute("merchant", market, routes)
	if err != nil || r.Providers[0] != "stripe" {
		t.Fatal("country fallback failed", r, err)
	}
	market.Currency = "USD"
	if _, err = SelectPaymentRoute("merchant", market, routes); err == nil {
		t.Fatal("unsupported currency selected")
	}
}
