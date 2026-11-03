package services

import (
	"elcomware/app/actions"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
)

type CatalogService struct {
	productQueries      *queries.ProductQueries
	createProductAction *actions.CreateProductAction
	deleteProductAction *actions.DeleteProductAction
}

func NewCatalogService() *CatalogService {
	return &CatalogService{
		productQueries:      queries.NewProductQueries(),
		createProductAction: actions.NewCreateProductAction(),
		deleteProductAction: actions.NewDeleteProductAction(),
	}
}

func (s *CatalogService) GetProducts(filter queries.ProductFilter) ([]models.Product, error) {
	return s.productQueries.FindFiltered(filter)
}

func (s *CatalogService) GetProductByID(id string) (*models.Product, error) {
	return s.productQueries.FindByID(id)
}

func (s *CatalogService) CreateProduct(req *requests.CreateProductRequest, user *models.User) (*models.Product, error) {
	return s.createProductAction.ExecuteFromRequest(req, user)
}

func (s *CatalogService) DeleteProduct(id string) error {
	return s.deleteProductAction.Execute(id)
}
