package actions

import (
	"elcomware/app/facades"
	"errors"
)

type DeleteProductAction struct{}

func NewDeleteProductAction() *DeleteProductAction {
	return &DeleteProductAction{}
}

func (a *DeleteProductAction) Execute(id string) error {
	result, err := facades.Orm().Query().Exec("DELETE FROM products WHERE id = ? AND NOT EXISTS (SELECT 1 FROM order_items WHERE product_id = ?) AND NOT EXISTS (SELECT 1 FROM offers WHERE product_id = ?)", id, id, id)
	if err != nil {
		return err
	}
	if result.RowsAffected != 1 {
		return errors.New("This listing has order or offer history; set stock to 0 instead.")
	}
	return nil
}
