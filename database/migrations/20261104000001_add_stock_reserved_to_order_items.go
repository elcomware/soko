package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261104000001AddStockReservedToOrderItems struct{}

func (m *M20261104000001AddStockReservedToOrderItems) Signature() string {
	return "20261104000001_add_stock_reserved_to_order_items"
}
func (m *M20261104000001AddStockReservedToOrderItems) Up() error {
	if facades.Schema().HasColumn("order_items", "stock_reserved") {
		return nil
	}
	return facades.Schema().Table("order_items", func(table schema.Blueprint) { table.Boolean("stock_reserved").Default(false) })
}
func (m *M20261104000001AddStockReservedToOrderItems) Down() error {
	if !facades.Schema().HasColumn("order_items", "stock_reserved") {
		return nil
	}
	return facades.Schema().DropColumns("order_items", []string{"stock_reserved"})
}
