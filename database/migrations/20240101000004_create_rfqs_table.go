package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"elcomware/app/facades"
)

type M20240101000004CreateRfqsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000004CreateRfqsTable) Signature() string {
	return "20240101000004_create_rfqs_table"
}

// Up Run the migrations.
func (r *M20240101000004CreateRfqsTable) Up() error {
	if !facades.Schema().HasTable("rfqs") {
		if err := facades.Schema().Create("rfqs", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("rfq_number", 64)
			table.String("buyer_id", 64).Nullable()
			table.String("buyer_name")
			table.String("buyer_company").Nullable()
			table.String("buyer_country", 64).Nullable()
			table.String("title")
			table.String("category", 128)
			table.Integer("quantity").Default(1)
			table.Double("target_price").Default(0)
			table.String("destination_country", 64).Nullable()
			table.LongText("specifications").Nullable()
			table.String("status", 32).Default("open")
			table.Timestamps()

			table.Index("rfq_number")
			table.Index("status")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("rfq_quotes") {
		if err := facades.Schema().Create("rfq_quotes", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("rfq_id", 64)
			table.String("vendor_id", 64).Nullable()
			table.String("vendor_name")
			table.String("vendor_country", 64).Nullable()
			table.Double("unit_price").Default(0)
			table.Double("total_amount").Default(0)
			table.Integer("lead_time_days").Default(7)
			table.LongText("notes").Nullable()
			table.String("status", 32).Default("pending")
			table.Timestamps()

			table.Index("rfq_id")
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20240101000004CreateRfqsTable) Down() error {
	if err := facades.Schema().DropIfExists("rfq_quotes"); err != nil {
		return err
	}
	return facades.Schema().DropIfExists("rfqs")
}
