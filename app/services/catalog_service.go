package services

import (
	"elcomware/app/actions"
	"elcomware/app/dto"
	"elcomware/app/mappers"
	"elcomware/app/queries"
)

type CatalogService struct {
	productQueries      *queries.ProductQueries
	productMapper       *mappers.ProductMapper
	createProductAction *actions.CreateProductAction
	deleteProductAction *actions.DeleteProductAction
}

func NewCatalogService() *CatalogService {
	return &CatalogService{
		productQueries:      queries.NewProductQueries(),
		productMapper:       mappers.NewProductMapper(),
		createProductAction: actions.NewCreateProductAction(),
		deleteProductAction: actions.NewDeleteProductAction(),
	}
}

func (s *CatalogService) GetCatalog(filter dto.ProductFilterInput) ([]dto.ProductOutput, error) {
	products, err := s.productQueries.FindFiltered(filter)
	if err != nil {
		return nil, err
	}
	return s.productMapper.ToDTOList(products), nil
}

func (s *CatalogService) GetProductByID(id string) (*dto.ProductOutput, error) {
	product, err := s.productQueries.FindByID(id)
	if err != nil {
		return nil, err
	}
	dto := s.productMapper.ToDTO(*product)
	return &dto, nil
}

func (s *CatalogService) CreateProduct(input dto.CreateProductInput) (dto.ProductOutput, error) {
	return s.createProductAction.Execute(input)
}

func (s *CatalogService) DeleteProduct(id string) error {
	return s.deleteProductAction.Execute(id)
}
