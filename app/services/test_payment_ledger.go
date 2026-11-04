package services

import (
	"context"
	"crypto/sha256"
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/payments"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"strings"
)

// RecordTestPaymentAttempt persists a server-created test session expectation.
// It must never be called with totals or identity supplied directly by a browser.
// No production order, stock, fulfillment or payout state is changed.
func RecordTestPaymentAttempt(a models.TestPaymentAttempt) error {
	if (a.ID == "" || len(a.ID) > 64) || a.Provider != "stripe_test" || (a.MerchantID == "" || len(a.MerchantID) > 64) || (a.OrderID == "" || len(a.OrderID) > 64) || (!strings.HasPrefix(a.PaymentReference, "cs_test_") || len(a.PaymentReference) > 255) || a.AmountMinor < 1 || len(a.Currency) != 3 {
		return errors.New("invalid isolated test attempt")
	}
	a.Status = "pending"
	a.Version = 0
	return facades.Orm().Query().Create(&a)
}

type TestPaymentLedger struct{}

// Handle verifies exact raw bytes before any durable write. It is intentionally
// not registered as an HTTP endpoint or activated via merchant preferences.
func (l TestPaymentLedger) Handle(ctx context.Context, adapter *payments.StripeTestAdapter, body []byte, headers map[string]string) (bool, error) {
	if adapter == nil {
		return false, payments.ErrUnconfigured
	}
	e, err := adapter.VerifyEvent(body, headers)
	if err != nil {
		return false, err
	}
	return l.ApplyVerified(ctx, e)
}

// ApplyVerified is an internal adapter boundary; never bind an Event from HTTP.
func (TestPaymentLedger) ApplyVerified(ctx context.Context, e payments.Event) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !e.TestMode || e.Provider != "stripe_test" || e.EventID == "" || len(e.EventID) > 512 || e.MerchantID == "" || e.OrderID == "" || e.AmountMinor < 1 {
		return false, errors.New("only verified isolated test events are accepted")
	}
	if e.Status != "pending" && e.Status != "succeeded" && e.Status != "failed" {
		return false, errors.New("unsupported payment state")
	}
	identity := sha256.Sum256([]byte(e.Provider + ":" + e.EventID))
	receiptID := hex.EncodeToString(identity[:])
	// Exclude arrival timestamp: reconciliation may report the same state later.
	canonical, _ := json.Marshal([]any{e.Provider, e.EventID, e.MerchantID, e.OrderID, e.PaymentReference, e.Currency, e.AmountMinor, e.Status, e.TestMode})
	digest := sha256.Sum256(canonical)
	eventHash := hex.EncodeToString(digest[:])
	duplicate := false
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		// Acquire the attempt's database write lock before read/transition/receipt.
		// This serializes across workers; unique receipts also bind an event globally.
		result, err := tx.Exec("UPDATE test_payment_attempts SET version = version + 1 WHERE payment_reference = ? AND provider = ?", e.PaymentReference, e.Provider)
		if err != nil {
			return err
		}
		if result.RowsAffected != 1 {
			return errors.New("unknown test payment attempt")
		}
		var attempt models.TestPaymentAttempt
		if err := tx.Where("payment_reference = ?", e.PaymentReference).First(&attempt); err != nil {
			return err
		}
		if attempt.MerchantID != e.MerchantID || attempt.OrderID != e.OrderID || attempt.Currency != e.Currency || attempt.AmountMinor != e.AmountMinor {
			return errors.New("test payment event does not match the expected order")
		}
		var receipt models.TestPaymentReceipt
		if err := tx.Where("id = ?", receiptID).First(&receipt); err != nil {
			return err
		}
		if receipt.ID != "" {
			if receipt.AttemptID != attempt.ID || receipt.EventHash != eventHash {
				return errors.New("event identity reused with different data")
			}
			duplicate = true
			return nil
		}
		// Paid is terminal; a delayed failure/pending event cannot undo success.
		next := attempt.Status
		if e.Status == "succeeded" || (attempt.Status == "pending" && e.Status == "failed") {
			next = e.Status
		}
		if _, err := tx.Model(&models.TestPaymentAttempt{}).Where("id = ?", attempt.ID).Update(map[string]any{"status": next}); err != nil {
			return err
		}
		return tx.Create(&models.TestPaymentReceipt{ID: receiptID, AttemptID: attempt.ID, EventHash: eventHash})
	})
	return duplicate, err
}
