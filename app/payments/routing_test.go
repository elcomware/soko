package payments

import (
	"context"
	"errors"
	"testing"
)

type testAdapter struct {
	Unconfigured
	id   string
	caps []Capability
	err  error
}

func (a testAdapter) ID() string { return a.id }
func (a testAdapter) Capabilities(context.Context, string) ([]Capability, error) {
	return a.caps, a.err
}
func TestRoutingRestrictsMerchantMarketAndVerifiedCapabilities(t *testing.T) {
	m := Market{"KE", "Nairobi", "KES", "mobile_money"}
	cap := Capability{Market: m, Collection: true}
	adapters := map[string]Adapter{"a": testAdapter{id: "a", caps: []Capability{cap}}, "b": testAdapter{id: "b", caps: []Capability{cap}}}
	route := Route{"merchant", m, []string{"b", "a"}}
	got, err := Candidates(context.Background(), "merchant", m, route, adapters)
	if err != nil || len(got) != 2 || got[0].Adapter.ID() != "b" {
		t.Fatalf("unexpected candidates %v %v", got, err)
	}
	if _, err = Candidates(context.Background(), "other", m, route, adapters); err == nil {
		t.Fatal("cross-merchant route")
	}
	for _, bad := range []Market{{"US", "Nairobi", "KES", "mobile_money"}, {"KE", "Mombasa", "KES", "mobile_money"}, {"KE", "Nairobi", "USD", "mobile_money"}, {"KE", "Nairobi", "KES", "card"}} {
		if _, err = Candidates(context.Background(), "merchant", bad, route, adapters); err == nil {
			t.Fatal("unsupported market accepted")
		}
	}
	adapters["b"] = testAdapter{id: "b", caps: []Capability{{Market: m, Collection: false}}}
	got, err = Candidates(context.Background(), "merchant", m, route, adapters)
	if err != nil || len(got) != 1 || got[0].Adapter.ID() != "a" {
		t.Fatal("inactive provider selected")
	}
	adapters["a"] = testAdapter{id: "a", err: errors.New("offline")}
	if _, err = Candidates(context.Background(), "merchant", m, route, adapters); err == nil {
		t.Fatal("capability lookup failed open")
	}
}
func TestRoutingMissingAdapterAndDuplicatesFailClosed(t *testing.T) {
	m := Market{"US", "", "USD", "card"}
	route := Route{"m", m, []string{"missing"}}
	if _, err := Candidates(context.Background(), "m", m, route, nil); !errors.Is(err, ErrUnconfigured) {
		t.Fatal(err)
	}
	a := testAdapter{id: "a", caps: []Capability{{Market: m, Collection: true}}}
	route.Providers = []string{"a", "a"}
	if _, err := Candidates(context.Background(), "m", m, route, map[string]Adapter{"a": a}); err == nil {
		t.Fatal("duplicate adapter accepted")
	}
}
