package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261104000020CreateProductMedia struct{}

func (m *M20261104000020CreateProductMedia) Signature() string {
	return "20261104000020_create_product_media"
}
func (m *M20261104000020CreateProductMedia) Up() error {
	if facades.Schema().HasTable("product_media") {
		return nil
	}
	return facades.Schema().Create("product_media", func(t schema.Blueprint) {
		t.String("id", 64)
		t.Primary("id")
		t.String("owner_id", 64)
		t.Index("owner_id")
		t.LongText("url")
		t.String("kind", 24)
		t.Integer("width")
		t.Integer("height")
		t.LongText("options_json")
		t.Timestamp("created_at")
	})
}
func (m *M20261104000020CreateProductMedia) Down() error {
	return facades.Schema().DropIfExists("product_media")
}
