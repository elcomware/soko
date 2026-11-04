package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261104000004CreatePromotions struct{}

func (m *M20261104000004CreatePromotions) Signature() string {
	return "20261104000004_create_promotions"
}
func (m *M20261104000004CreatePromotions) Up() error {
	if !facades.Schema().HasTable("promotions") {
		if err := facades.Schema().Create("promotions", func(t schema.Blueprint) {
			t.String("id", 64)
			t.Primary("id")
			t.String("name", 120)
			t.String("vendor_id", 64).Default("")
			t.String("kind", 16)
			t.Double("amount")
			t.Boolean("enabled").Default(false)
			t.Timestamp("starts_at")
			t.Timestamp("ends_at")
			t.LongText("scope_json")
			t.String("created_by", 64)
			t.Timestamps()
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasColumn("order_items", "original_unit_price") {
		return facades.Schema().Table("order_items", func(t schema.Blueprint) {
			t.Double("original_unit_price").Default(0)
			t.Double("discount_per_unit").Default(0)
			t.String("promotion_id", 64).Default("")
			t.String("promotion_name", 120).Default("")
		})
	}
	return nil
}
func (m *M20261104000004CreatePromotions) Down() error {
	if facades.Schema().HasColumn("order_items", "original_unit_price") {
		if err := facades.Schema().DropColumns("order_items", []string{"original_unit_price", "discount_per_unit", "promotion_id", "promotion_name"}); err != nil {
			return err
		}
	}
	return facades.Schema().DropIfExists("promotions")
}
