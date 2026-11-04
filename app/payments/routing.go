package payments

import (
	"context"
	"errors"
	"fmt"
)

// Market is supplied by trusted checkout resolution, not a browser routing hint.
// Area is a provider-defined subdivision code; empty means country-wide.
type Market struct{ Country, Area, Currency, Method string }
type Capability struct {
	Market          Market
	Collection      bool
	Refunds         bool
	SplitSettlement bool
}

// Adapter isolates gateway SDKs and exposes only account-verified capabilities.
// Capability combinations are explicit: methods/currencies must not be expanded
// into a Cartesian product that a provider may not actually support.
type Adapter interface {
	Provider
	ID() string
	Capabilities(context.Context, string) ([]Capability, error)
}
type Route struct {
	MerchantID string
	Market     Market
	Providers  []string // Merchant-approved preference order; no implicit fallback.
}
type Candidate struct {
	Adapter    Adapter
	Capability Capability
}

var ErrNoRoute = errors.New("no configured payment route supports this merchant and market")

// Candidates resolves an ordered allowlist; it neither collects payment nor
// retries another gateway. An ambiguous attempt must be reconciled with the
// original provider before any new attempt, to prevent double collection.
func Candidates(ctx context.Context, merchantID string, market Market, route Route, adapters map[string]Adapter) ([]Candidate, error) {
	if merchantID == "" || route.MerchantID != merchantID || market.Country == "" || market.Currency == "" || market.Method == "" || !matches(route.Market, market) {
		return nil, ErrNoRoute
	}
	result := []Candidate{}
	seen := map[string]bool{}
	for _, id := range route.Providers {
		if id == "" || seen[id] {
			return nil, errors.New("payment route contains duplicate or empty provider")
		}
		seen[id] = true
		adapter, ok := adapters[id]
		if !ok || adapter == nil || adapter.ID() != id {
			return nil, ErrUnconfigured
		}
		capabilities, err := adapter.Capabilities(ctx, merchantID)
		if err != nil {
			return nil, fmt.Errorf("provider %s capabilities unavailable: %w", id, err)
		}
		for _, capability := range capabilities {
			if capability.Collection && matches(capability.Market, market) {
				result = append(result, Candidate{adapter, capability})
				break
			}
		}
	}
	if len(result) == 0 {
		return nil, ErrNoRoute
	}
	return result, nil
}
func matches(supported, requested Market) bool {
	return supported.Country == requested.Country && supported.Currency == requested.Currency && supported.Method == requested.Method && (supported.Area == "" || supported.Area == requested.Area)
}
