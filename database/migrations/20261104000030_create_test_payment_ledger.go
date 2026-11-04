package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261104000030CreateTestPaymentLedger struct{}

func (*M20261104000030CreateTestPaymentLedger) Signature() string {
	return "20261104000030_create_test_payment_ledger"
}
func (*M20261104000030CreateTestPaymentLedger) Up() error {
	if !facades.Schema().HasTable("test_payment_attempts") {
		if err := facades.Schema().Create("test_payment_attempts", func(t schema.Blueprint) {
			t.String("id", 64)
			t.Primary("id")
			t.String("provider", 32)
			t.String("merchant_id", 64)
			t.String("order_id", 64)
			t.String("payment_reference", 255)
			t.Unique("payment_reference")
			t.String("currency", 3)
			t.BigInteger("amount_minor")
			t.String("status", 32)
			t.Integer("version").Default(0)
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasTable("test_payment_receipts") {
		return facades.Schema().Create("test_payment_receipts", func(t schema.Blueprint) {
			t.String("id", 64)
			t.Primary("id")
			t.String("attempt_id", 64)
			t.String("event_hash", 64)
		})
	}
	return nil
}
func (*M20261104000030CreateTestPaymentLedger) Down() error {
	if err := facades.Schema().DropIfExists("test_payment_receipts"); err != nil {
		return err
	}
	return facades.Schema().DropIfExists("test_payment_attempts")
}
