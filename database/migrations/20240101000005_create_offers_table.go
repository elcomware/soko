package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"elcomware/app/facades"
)

type M20240101000005CreateOffersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000005CreateOffersTable) Signature() string {
	return "20240101000005_create_offers_table"
}

// Up Run the migrations.
func (r *M20240101000005CreateOffersTable) Up() error {
	if !facades.Schema().HasTable("offers") {
		return facades.Schema().Create("offers", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("product_id", 64)
			table.String("product_title")
			table.String("buyer_id", 64).Nullable()
			table.String("buyer_name")
			table.String("buyer_email").Nullable()
			table.String("buyer_company").Nullable()
			table.Double("offered_price").Default(0)
			table.Integer("quantity").Default(1)
			table.Double("total_offered").Default(0)
			table.LongText("message").Nullable()
			table.String("status", 32).Default("pending")
			table.Timestamps()

			table.Index("product_id")
			table.Index("status")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20240101000005CreateOffersTable) Down() error {
	return facades.Schema().DropIfExists("offers")
}

