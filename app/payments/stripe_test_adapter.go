package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StripeTestAdapter implements the real Checkout HTTP protocol for isolated test
// accounts only. It is never registered with the application's live route resolver.
// Credentials are supplied by server-side dependency injection, never forms.
type StripeTestAdapter struct {
	key, webhookSecret, merchant, currency, successURL, cancelURL string
	client                                                        *http.Client
	now                                                           func() time.Time
}
type StripeTestConfig struct{ SecretKey, WebhookSecret, MerchantID, Currency, SuccessURL, CancelURL string }

func NewStripeTestAdapter(c StripeTestConfig, transport http.RoundTripper) (*StripeTestAdapter, error) {
	if !strings.HasPrefix(c.SecretKey, "sk_test_") || !strings.HasPrefix(c.WebhookSecret, "whsec_") || (c.MerchantID == "" || len(c.MerchantID) > 64) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(c.Currency) {
		return nil, ErrUnconfigured
	}
	for _, raw := range []string{c.SuccessURL, c.CancelURL} {
		u, e := url.Parse(raw)
		if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
			return nil, errors.New("test return URLs must use HTTPS")
		}
	}
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &StripeTestAdapter{key: c.SecretKey, webhookSecret: c.WebhookSecret, merchant: c.MerchantID, currency: c.Currency, successURL: c.SuccessURL, cancelURL: c.CancelURL, client: &http.Client{Transport: transport, Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("payment API redirects are forbidden") }}, now: time.Now}, nil
}
func (*StripeTestAdapter) ID() string { return "stripe_test" }
func (*StripeTestAdapter) Capabilities(context.Context, string) ([]Capability, error) {
	return nil, ErrUnconfigured
} // Never eligible for storefront collection.
type stripeSession struct {
	ID            string            `json:"id"`
	URL           string            `json:"url"`
	LiveMode      *bool             `json:"livemode"`
	Amount        int64             `json:"amount_total"`
	Currency      string            `json:"currency"`
	PaymentStatus string            `json:"payment_status"`
	Status        string            `json:"status"`
	Mode          string            `json:"mode"`
	Metadata      map[string]string `json:"metadata"`
}

func (s *StripeTestAdapter) request(ctx context.Context, method, path string, values url.Values, key string) (stripeSession, error) {
	var session stripeSession
	req, err := http.NewRequestWithContext(ctx, method, "https://api.stripe.com/v1/checkout/sessions"+path, strings.NewReader(values.Encode()))
	if err != nil {
		return session, err
	}
	req.SetBasicAuth(s.key, "")
	req.Header.Set("Stripe-Version", "2026-10-28.endive")
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Idempotency-Key", key)
	}
	response, err := s.client.Do(req)
	if err != nil {
		return session, errors.New("Stripe test request outcome unknown; reconcile before retrying")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024+1))
	if err != nil || len(body) > 1024*1024 {
		return session, errors.New("invalid Stripe test response")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return session, fmt.Errorf("Stripe test request returned HTTP %d; reconcile ambiguous outcomes", response.StatusCode)
	}
	if json.Unmarshal(body, &session) != nil {
		return session, errors.New("malformed Stripe test session")
	}
	if err = s.validateSession(session); err != nil {
		return stripeSession{}, err
	}
	return session, nil
}
func (s *StripeTestAdapter) validateSession(session stripeSession) error {
	if session.LiveMode == nil || *session.LiveMode || (!regexp.MustCompile(`^cs_test_[A-Za-z0-9_]+$`).MatchString(session.ID) || len(session.ID) > 255) || session.Mode != "payment" || session.Metadata["merchant_id"] != s.merchant || (session.Metadata["order_id"] == "" || len(session.Metadata["order_id"]) > 64) || strings.ToUpper(session.Currency) != s.currency || session.Amount < 1 {
		return errors.New("Stripe session failed test-account identity checks")
	}
	return nil
}
func (s *StripeTestAdapter) CreateSession(ctx context.Context, r Request) (Session, error) {
	if r.MerchantID != s.merchant || r.Currency != s.currency || r.Method != "card" || r.AmountMinor < 1 || r.AmountMinor > 99999999 || r.OrderID == "" || len(r.OrderID) > 64 || r.IdempotencyKey == "" || len(r.IdempotencyKey) > 255 {
		return Session{}, errors.New("invalid isolated test checkout request")
	}
	values := url.Values{"mode": {"payment"}, "success_url": {s.successURL}, "cancel_url": {s.cancelURL}, "payment_method_types[0]": {"card"}, "client_reference_id": {r.OrderID}, "metadata[merchant_id]": {r.MerchantID}, "metadata[order_id]": {r.OrderID}, "line_items[0][price_data][currency]": {strings.ToLower(r.Currency)}, "line_items[0][price_data][unit_amount]": {strconv.FormatInt(r.AmountMinor, 10)}, "line_items[0][price_data][product_data][name]": {"Soko test order"}, "line_items[0][quantity]": {"1"}}
	session, err := s.request(ctx, http.MethodPost, "", values, r.IdempotencyKey)
	if err != nil {
		return Session{}, err
	}
	if session.Amount != r.AmountMinor || session.Metadata["order_id"] != r.OrderID {
		return Session{}, errors.New("Stripe test quote mismatch")
	}
	u, err := url.Parse(session.URL)
	if err != nil || u.Scheme != "https" || u.Host != "checkout.stripe.com" || u.User != nil {
		return Session{}, errors.New("invalid Stripe hosted checkout URL")
	}
	return Session{ProviderReference: session.ID, RedirectURL: session.URL}, nil
}
func (s *StripeTestAdapter) event(session stripeSession, id string, at time.Time) (Event, error) {
	if err := s.validateSession(session); err != nil {
		return Event{}, err
	}
	status := "pending"
	if session.PaymentStatus == "paid" {
		status = "succeeded"
	} else if session.Status == "expired" {
		status = "failed"
	}
	return Event{Provider: s.ID(), EventID: id, MerchantID: s.merchant, OrderID: session.Metadata["order_id"], PaymentReference: session.ID, Currency: s.currency, AmountMinor: session.Amount, Status: status, OccurredAt: at, TestMode: true}, nil
}
func (s *StripeTestAdapter) VerifyEvent(body []byte, headers map[string]string) (Event, error) {
	signature := ""
	for key, value := range headers {
		if strings.EqualFold(key, "Stripe-Signature") {
			signature = value
		}
	}
	if err := VerifyStripeSignature(body, signature, s.webhookSecret, s.now()); err != nil {
		return Event{}, err
	}
	var envelope struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		LiveMode *bool  `json:"livemode"`
		Created  int64  `json:"created"`
		Data     struct {
			Object stripeSession `json:"object"`
		} `json:"data"`
	}
	if json.Unmarshal(body, &envelope) != nil || envelope.ID == "" || len(envelope.ID) > 255 || envelope.LiveMode == nil || *envelope.LiveMode {
		return Event{}, errors.New("invalid test webhook envelope")
	}
	switch envelope.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded", "checkout.session.async_payment_failed", "checkout.session.expired":
	default:
		return Event{}, errors.New("unsupported test webhook type")
	}
	event, err := s.event(envelope.Data.Object, envelope.ID, time.Unix(envelope.Created, 0))
	if err != nil {
		return Event{}, err
	}
	if envelope.Type == "checkout.session.async_payment_failed" && event.Status != "succeeded" {
		event.Status = "failed"
	}
	return event, nil
}
func (s *StripeTestAdapter) Reconcile(ctx context.Context, reference string) (Event, error) {
	if !regexp.MustCompile(`^cs_test_[A-Za-z0-9_]+$`).MatchString(reference) {
		return Event{}, errors.New("invalid test session reference")
	}
	session, err := s.request(ctx, http.MethodGet, "/"+reference, nil, "")
	if err != nil {
		return Event{}, err
	}
	if session.ID != reference {
		return Event{}, errors.New("reconciliation reference mismatch")
	}
	return s.event(session, "reconcile:"+session.ID+":"+session.PaymentStatus+":"+session.Status, s.now())
}
