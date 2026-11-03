package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"elcomware/app/facades"
)

type M20240101000003CreateOrdersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000003CreateOrdersTable) Signature() string {
	return "20240101000003_create_orders_table"
}

// Up Run the migrations.
func (r *M20240101000003CreateOrdersTable) Up() error {
	if !facades.Schema().HasTable("orders") {
		if err := facades.Schema().Create("orders", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("order_number", 64)
			table.String("user_id", 64).Nullable()
			table.String("customer_name")
			table.String("customer_email")
			table.Double("total_amount").Default(0)
			table.String("currency", 16).Default("USD")
			table.String("status", 32).Default("pending")
			table.String("payment_method", 64).Default("Trade Assurance Escrow")
			table.Boolean("has_digital").Default(false)
			table.Boolean("has_physical").Default(false)
			table.Timestamps()

			table.Index("order_number")
			table.Index("user_id")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("order_items") {
		if err := facades.Schema().Create("order_items", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("order_id", 64)
			table.String("product_id", 64)
			table.String("product_title")
			table.String("product_type", 64).Default("physical")
			table.Integer("quantity").Default(1)
			table.Double("unit_price").Default(0)
			table.Double("total_price").Default(0)
			table.Text("digital_license_key").Nullable()
			table.Text("digital_download_url").Nullable()
			table.Timestamps()

			table.Index("order_id")
			table.Index("product_id")
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20240101000003CreateOrdersTable) Down() error {
	if err := facades.Schema().DropIfExists("order_items"); err != nil {
		return err
	}
	return facades.Schema().DropIfExists("orders")
}
