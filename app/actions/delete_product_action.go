package actions

import (
	"elcomware/app/facades"
	"elcomware/app/models"
)

type DeleteProductAction struct{}

func NewDeleteProductAction() *DeleteProductAction {
	return &DeleteProductAction{}
}

func (a *DeleteProductAction) Execute(id string) error {
	_, err := facades.Orm().Query().Where("id = ?", id).Delete(&models.Product{})
	return err
}
