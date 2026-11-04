// Package payments defines the server-side boundary for optional payment engines.
// Preferences alone never enable collection. An adapter must verify merchant
// onboarding, supported country/currency/methods and signed events before use.
package payments

import (
	"context"
	"errors"
	"time"
)

var ErrUnconfigured = errors.New("payment provider is not configured")

type Request struct {
	MerchantID, OrderID, Currency, Method, IdempotencyKey string
	AmountMinor                                           int64
}
type Session struct{ ProviderReference, RedirectURL string }
type Event struct {
	Provider, EventID, MerchantID, OrderID, PaymentReference, Currency, Status string
	AmountMinor                                                                int64
	OccurredAt                                                                 time.Time
	TestMode                                                                   bool
}
type Provider interface {
	CreateSession(context.Context, Request) (Session, error)
	VerifyEvent([]byte, map[string]string) (Event, error)
	Reconcile(context.Context, string) (Event, error)
}

// EventStore implementations must persist a unique (provider,event_id) receipt
// atomically with order transitions, rejecting wrong amount/currency/merchant.
type EventStore interface {
	ApplyVerified(context.Context, Event) (alreadyProcessed bool, err error)
}
type Unconfigured struct{}

func (Unconfigured) CreateSession(context.Context, Request) (Session, error) {
	return Session{}, ErrUnconfigured
}
func (Unconfigured) VerifyEvent([]byte, map[string]string) (Event, error) {
	return Event{}, ErrUnconfigured
}
func (Unconfigured) Reconcile(context.Context, string) (Event, error) {
	return Event{}, ErrUnconfigured
}
