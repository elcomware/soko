package actions

import (
	"fmt"
	"strings"
	"time"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
)

type CreateProductAction struct {
	mapper *mappers.ProductMapper
}

func NewCreateProductAction() *CreateProductAction {
	return &CreateProductAction{
		mapper: mappers.NewProductMapper(),
	}
}

func (a *CreateProductAction) Execute(input dto.CreateProductInput) (dto.ProductOutput, error) {
	product := a.mapper.ToModel(input)
	product.ID = fmt.Sprintf("prod_%d", time.Now().UnixNano())
	if product.Slug == "" {
		slug := strings.ToLower(strings.ReplaceAll(product.Title, " ", "-"))
		product.Slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix()%10000)
	}

	if err := facades.Orm().Query().Create(&product); err != nil {
		return dto.ProductOutput{}, err
	}

	return a.mapper.ToDTO(product), nil
}
