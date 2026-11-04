package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261104000031CreateMerchandising struct{}

func (*M20261104000031CreateMerchandising) Signature() string {
	return "20261104000031_create_merchandising"
}
func (*M20261104000031CreateMerchandising) Up() error {
	if facades.Schema().HasTable("storefront_merchandising") {
		return nil
	}
	return facades.Schema().Create("storefront_merchandising", func(t schema.Blueprint) {
		t.String("id", 32)
		t.Primary("id")
		t.LongText("settings_json")
		t.Integer("version").Default(1)
	})
}
func (*M20261104000031CreateMerchandising) Down() error {
	return facades.Schema().DropIfExists("storefront_merchandising")
}
