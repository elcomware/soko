package payments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

type mockStripeTransport func(*http.Request) (*http.Response, error)

func (f mockStripeTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func testStripeConfig() StripeTestConfig {
	return StripeTestConfig{SecretKey: "sk_test_fixture_not_a_credential", WebhookSecret: "whsec_fixture_not_a_credential", MerchantID: "merchant", Currency: "USD", SuccessURL: "https://shop.example/test-success", CancelURL: "https://shop.example/test-cancel"}
}
func testSessionJSON() string {
	return `{"id":"cs_test_fixture","url":"https://checkout.stripe.com/c/pay/cs_test_fixture","livemode":false,"amount_total":1500,"currency":"usd","mode":"payment","status":"complete","payment_status":"paid","metadata":{"merchant_id":"merchant","order_id":"order"}}`
}
func signedTestEvent(t *testing.T, a *StripeTestAdapter, session, kind string) ([]byte, map[string]string) {
	t.Helper()
	body := []byte(`{"id":"evt_fixture","livemode":false,"created":1700000000,"type":"` + kind + `","data":{"object":` + session + `}}`)
	at := strconv.FormatInt(a.now().Unix(), 10)
	h := hmac.New(sha256.New, []byte(a.webhookSecret))
	h.Write([]byte(at + "."))
	h.Write(body)
	return body, map[string]string{"Stripe-Signature": "t=" + at + ",v1=" + hex.EncodeToString(h.Sum(nil))}
}
func TestStripeTestAdapterCreateAndReconcile(t *testing.T) {
	calls := 0
	a, err := NewStripeTestAdapter(testStripeConfig(), mockStripeTransport(func(req *http.Request) (*http.Response, error) {
		calls++
		if req.URL.Host != "api.stripe.com" || req.Header.Get("Stripe-Version") == "" {
			t.Fatal("unexpected endpoint")
		}
		key, _, ok := req.BasicAuth()
		if !ok || key != testStripeConfig().SecretKey {
			t.Fatal("missing auth")
		}
		if req.Method == "POST" {
			if req.Header.Get("Idempotency-Key") != "stable" {
				t.Fatal("missing stable key")
			}
			if err := req.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if req.Form.Get("line_items[0][price_data][unit_amount]") != "1500" || req.Form.Get("metadata[merchant_id]") != "merchant" {
				t.Fatal("wrong request fields")
			}
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(testSessionJSON())), Header: make(http.Header)}, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	r := Request{MerchantID: "merchant", OrderID: "order", Currency: "USD", Method: "card", IdempotencyKey: "stable", AmountMinor: 1500}
	session, err := a.CreateSession(context.Background(), r)
	if err != nil || session.ProviderReference != "cs_test_fixture" {
		t.Fatal(session, err)
	}
	e, err := a.Reconcile(context.Background(), session.ProviderReference)
	if err != nil || !e.TestMode || e.Status != "succeeded" {
		t.Fatal(e, err)
	}
	r.OrderID = strings.Repeat("o", 65)
	if _, err := a.CreateSession(context.Background(), r); err == nil {
		t.Fatal("overlong order ID accepted")
	}
	r.OrderID = "order"
	r.MerchantID = "other"
	if _, err := a.CreateSession(context.Background(), r); err == nil {
		t.Fatal("cross merchant accepted")
	}
	if calls != 2 {
		t.Fatal("invalid request reached API")
	}
	if _, err := a.Capabilities(context.Background(), "merchant"); err == nil {
		t.Fatal("test adapter exposed live capabilities")
	}
}
func TestStripeTestAdapterRejectsLiveKeysAndUntrustedResponses(t *testing.T) {
	c := testStripeConfig()
	c.SecretKey = "sk_live_fixture"
	if _, err := NewStripeTestAdapter(c, nil); err == nil {
		t.Fatal("live key accepted")
	}
	for _, response := range []string{strings.Replace(testSessionJSON(), `"livemode":false`, `"livemode":true`, 1), strings.Replace(testSessionJSON(), "checkout.stripe.com", "evil.example", 1), strings.Replace(testSessionJSON(), "1500", "1600", 1), strings.Replace(testSessionJSON(), `"merchant"`, `"other"`, 1)} {
		a, _ := NewStripeTestAdapter(testStripeConfig(), mockStripeTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(response))}, nil
		}))
		if _, err := a.CreateSession(context.Background(), Request{MerchantID: "merchant", OrderID: "order", Currency: "USD", Method: "card", IdempotencyKey: "stable", AmountMinor: 1500}); err == nil {
			t.Fatal("bad response accepted", response)
		}
	}
}
func TestStripeTestWebhookRequiresPaidStateAndSignature(t *testing.T) {
	a, _ := NewStripeTestAdapter(testStripeConfig(), nil)
	a.now = func() time.Time { return time.Unix(1700000000, 0) }
	body, headers := signedTestEvent(t, a, testSessionJSON(), "checkout.session.completed")
	e, err := a.VerifyEvent(body, headers)
	if err != nil || e.Status != "succeeded" || !e.TestMode {
		t.Fatal(e, err)
	}
	if _, err = a.VerifyEvent(append(body, ' '), headers); err == nil {
		t.Fatal("tampered event accepted")
	}
	unpaid := strings.Replace(testSessionJSON(), `"payment_status":"paid"`, `"payment_status":"unpaid"`, 1)
	body, headers = signedTestEvent(t, a, unpaid, "checkout.session.completed")
	e, err = a.VerifyEvent(body, headers)
	if err != nil || e.Status != "pending" {
		t.Fatal("completed alone marked paid", e, err)
	}
	var obj map[string]any
	_ = json.Unmarshal([]byte(testSessionJSON()), &obj)
	obj["livemode"] = true
	raw, _ := json.Marshal(obj)
	body, headers = signedTestEvent(t, a, string(raw), "checkout.session.completed")
	if _, err = a.VerifyEvent(body, headers); err == nil {
		t.Fatal("live event accepted")
	}
}
