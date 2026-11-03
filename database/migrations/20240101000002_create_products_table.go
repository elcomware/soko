package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"elcomware/app/facades"
)

type M20240101000002CreateProductsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000002CreateProductsTable) Signature() string {
	return "20240101000002_create_products_table"
}

// Up Run the migrations.
func (r *M20240101000002CreateProductsTable) Up() error {
	if !facades.Schema().HasTable("products") {
		return facades.Schema().Create("products", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("title")
			table.String("slug").Nullable()
			table.LongText("description").Nullable()
			table.Text("short_desc").Nullable()
			table.String("category", 128)
			table.String("sub_category", 128).Nullable()
			table.String("type", 64).Default("physical")
			table.String("format", 64).Default("wholesale")
			table.String("condition", 64).Default("new")
			table.Double("price").Default(0)
			table.Double("wholesale_price").Default(0)
			table.LongText("price_tiers").Nullable()
			table.Integer("moq").Default(1)
			table.Integer("stock").Default(0)
			table.String("sku", 64).Nullable()
			table.String("vendor_id", 64).Nullable()
			table.String("vendor_name", 128).Nullable()
			table.String("vendor_country", 64).Nullable()
			table.String("vendor_region", 64).Nullable()
			table.Boolean("vendor_verified").Default(false)
			table.String("vendor_badge", 64).Nullable()
			table.String("vendor_positive_rate", 32).Nullable()
			table.Double("trade_assurance_cover").Default(0)
			table.Boolean("accepts_offers").Default(false)
			table.Double("rating").Default(5.0)
			table.Integer("reviews_count").Default(0)
			table.Text("image_url").Nullable()
			table.LongText("gallery_urls").Nullable()
			table.Text("tags").Nullable()
			table.Boolean("featured").Default(false)
			table.Integer("lead_time_days").Default(0)
			table.String("incoterm", 32).Nullable()
			table.Boolean("local_pickup_available").Default(false)
			table.Text("digital_download_url").Nullable()
			table.String("digital_license_format", 64).Nullable()
			table.String("digital_file_size", 64).Nullable()
			table.Timestamps()

			table.Index("category")
			table.Index("type")
			table.Index("format")
			table.Index("vendor_region")
			table.Index("featured")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20240101000002CreateProductsTable) Down() error {
	return facades.Schema().DropIfExists("products")
}

