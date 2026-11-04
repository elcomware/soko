package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261103000001CreateStorefrontContents struct{}

func (m *M20261103000001CreateStorefrontContents) Signature() string {
	return "20261103000001_create_storefront_contents"
}
func (m *M20261103000001CreateStorefrontContents) Up() error {
	if facades.Schema().HasTable("storefront_contents") {
		return nil
	}
	return facades.Schema().Create("storefront_contents", func(table schema.Blueprint) {
		table.String("id", 32)
		table.Primary("id")
		table.LongText("content_json")
		table.Timestamps()
	})
}
func (m *M20261103000001CreateStorefrontContents) Down() error {
	return facades.Schema().DropIfExists("storefront_contents")
}
