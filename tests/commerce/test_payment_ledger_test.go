package commerce_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/payments"
	"elcomware/app/services"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestVerifiedTestPaymentLedgerDeduplicatesAndProtectsOrder(t *testing.T) {
	attempt := models.TestPaymentAttempt{ID: "sandbox-attempt", Provider: "stripe_test", MerchantID: "sandbox-merchant", OrderID: "sandbox-order", PaymentReference: "cs_test_ledger", Currency: "USD", AmountMinor: 2500}
	require.NoError(t, services.RecordTestPaymentAttempt(attempt))
	order := models.Order{ID: attempt.OrderID, OrderNumber: "sandbox-order-number", UserID: "sandbox-buyer", Currency: "USD", TotalAmount: 25, Status: models.OrderStatusPending}
	require.NoError(t, facades.Orm().Query().Create(&order))
	a, err := payments.NewStripeTestAdapter(payments.StripeTestConfig{SecretKey: "sk_test_fixture_not_a_credential", WebhookSecret: "whsec_fixture_not_a_credential", MerchantID: attempt.MerchantID, Currency: "USD", SuccessURL: "https://shop.example/success", CancelURL: "https://shop.example/cancel"}, nil)
	require.NoError(t, err)
	body := []byte(fmt.Sprintf(`{"id":"evt_ledger","type":"checkout.session.completed","livemode":false,"created":%d,"data":{"object":{"id":"cs_test_ledger","livemode":false,"amount_total":2500,"currency":"usd","mode":"payment","payment_status":"paid","status":"complete","metadata":{"merchant_id":"sandbox-merchant","order_id":"sandbox-order"}}}}`, time.Now().Unix()))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte("whsec_fixture_not_a_credential"))
	mac.Write([]byte(timestamp + "."))
	mac.Write(body)
	headers := map[string]string{"Stripe-Signature": "t=" + timestamp + ",v1=" + hex.EncodeToString(mac.Sum(nil))}
	ledger := services.TestPaymentLedger{}
	duplicate, err := ledger.Handle(context.Background(), a, body, headers)
	require.NoError(t, err)
	require.False(t, duplicate)
	duplicate, err = ledger.Handle(context.Background(), a, body, headers)
	require.NoError(t, err)
	require.True(t, duplicate)
	var saved models.TestPaymentAttempt
	require.NoError(t, facades.Orm().Query().Where("id = ?", attempt.ID).First(&saved))
	require.Equal(t, "succeeded", saved.Status)
	var unchanged models.Order
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.ID).First(&unchanged))
	require.Equal(t, models.OrderStatusPending, unchanged.Status)
	receipts, countErr := facades.Orm().Query().Model(&models.TestPaymentReceipt{}).Where("attempt_id = ?", attempt.ID).Count()
	require.NoError(t, countErr)
	require.EqualValues(t, 1, receipts)
	_, err = ledger.Handle(context.Background(), a, append(body, ' '), headers)
	require.Error(t, err)
}
func TestTestPaymentLedgerRejectsMismatchAndLateFailureCannotUndoPaid(t *testing.T) {
	a := models.TestPaymentAttempt{ID: "sandbox-mismatch", Provider: "stripe_test", MerchantID: "merchant", OrderID: "order", PaymentReference: "cs_test_mismatch", Currency: "USD", AmountMinor: 1000}
	require.NoError(t, services.RecordTestPaymentAttempt(a))
	ledger := services.TestPaymentLedger{}
	e := payments.Event{Provider: "stripe_test", TestMode: true, EventID: "evt_mismatch", MerchantID: a.MerchantID, OrderID: a.OrderID, PaymentReference: a.PaymentReference, Currency: a.Currency, AmountMinor: a.AmountMinor, Status: "succeeded"}
	for _, change := range []func(*payments.Event){func(e *payments.Event) { e.AmountMinor++ }, func(e *payments.Event) { e.MerchantID = "other" }, func(e *payments.Event) { e.OrderID = "other" }, func(e *payments.Event) { e.Currency = "EUR" }, func(e *payments.Event) { e.TestMode = false }} {
		bad := e
		change(&bad)
		_, err := ledger.ApplyVerified(context.Background(), bad)
		require.Error(t, err)
	}
	var saved models.TestPaymentAttempt
	require.NoError(t, facades.Orm().Query().Where("id = ?", a.ID).First(&saved))
	require.Equal(t, "pending", saved.Status)
	require.Equal(t, 0, saved.Version)
	_, err := ledger.ApplyVerified(context.Background(), e)
	require.NoError(t, err)
	altered := e
	altered.Status = "failed"
	_, err = ledger.ApplyVerified(context.Background(), altered)
	require.Error(t, err, "event ID cannot be reused")
	late := e
	late.EventID = "evt_late_failure"
	late.Status = "failed"
	_, err = ledger.ApplyVerified(context.Background(), late)
	require.NoError(t, err)
	require.NoError(t, facades.Orm().Query().Where("id = ?", a.ID).First(&saved))
	require.Equal(t, "succeeded", saved.Status)
}
