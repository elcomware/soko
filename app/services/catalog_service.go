package services

import (
	"elcomware/app/actions"
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
	"errors"
	"net/url"
	"strings"
	"time"
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
	products, err := s.productQueries.FindFiltered(filter)
	if err == nil {
		err = queries.ApplyProductPricing(products)
	}
	return products, err
}

func (s *CatalogService) GetProductByID(id string) (*models.Product, error) {
	product, err := s.productQueries.FindByID(id)
	if err != nil {
		return nil, err
	}
	list := []models.Product{*product}
	if err = queries.ApplyProductPricing(list); err != nil {
		return nil, err
	}
	return &list[0], nil
}

func (s *CatalogService) CreateProduct(req *requests.CreateProductRequest, user *models.User) (*models.Product, error) {
	if err := validateProductMediaPolicy(req, nil, user); err != nil {
		return nil, err
	}
	return s.createProductAction.ExecuteFromRequest(req, user)
}

func (s *CatalogService) DeleteProduct(id string, user *models.User) error {
	p, err := s.GetProductByID(id)
	if err != nil {
		return err
	}
	if p == nil || !user.Can("products.delete", p.VendorID) {
		return errors.New("listing not found or access denied")
	}
	return s.deleteProductAction.Execute(id)
}
func (s *CatalogService) UpdateProduct(id string, r *requests.CreateProductRequest, user *models.User) (*models.Product, error) {
	p, err := s.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	if !user.CanManageProduct(p) {
		return nil, errors.New("listing not found or access denied")
	}
	if err := validateProductMediaPolicy(r, p, user); err != nil {
		return nil, err
	}
	if err := actions.ValidateProduct(r); err != nil {
		return nil, err
	}
	version, err := time.Parse(time.RFC3339Nano, r.Version)
	if err != nil || !version.Equal(p.UpdatedAt) {
		return nil, errors.New("This listing changed. Refresh it before saving your edits.")
	}
	previousUpdatedAt, previousStock := p.UpdatedAt, p.Stock
	actions.ApplyProductInput(p, r)
	values := map[string]any{"title": p.Title, "category": p.Category, "sub_category": p.SubCategory, "description": p.Description, "short_desc": p.ShortDesc, "price": p.Price, "wholesale_price": p.WholesalePrice, "moq": p.MOQ, "stock": p.Stock, "sku": p.Sku, "image_url": p.ImageUrl, "gallery_urls": p.GalleryUrlsJSON, "featured": p.Featured, "type": p.Type, "format": p.Format, "condition": p.Condition, "vendor_region": p.VendorRegion, "accepts_offers": p.AcceptsOffers, "digital_download_url": p.DigitalDownloadUrl, "digital_license_format": p.DigitalLicenseFormat, "digital_file_size": p.DigitalFileSize, "price_tiers": p.PriceTiersJSON, "lead_time_days": p.LeadTimeDays, "incoterm": p.Incoterm, "updated_at": p.UpdatedAt}
	result, err := facades.Orm().Query().Model(&models.Product{}).Where("id = ? AND updated_at = ? AND stock = ?", id, previousUpdatedAt, previousStock).Update(values)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected != 1 {
		return nil, errors.New("This listing changed. Refresh it before saving your edits.")
	}
	return p, nil
}

func validateProductMediaPolicy(r *requests.CreateProductRequest, previous *models.Product, actor *models.User) error {
	p := LoadMediaPolicy()
	for _, raw := range append([]string{r.ImageUrl}, r.GalleryUrls...) {
		if err := ValidateProductImageURL(raw); err != nil {
			return err
		}
		existing := previous != nil && raw == previous.ImageUrl
		if previous != nil {
			previous.Hydrate()
			for _, old := range previous.GalleryUrls {
				if old == raw {
					existing = true
				}
			}
		}
		parsed, _ := url.Parse(raw)
		if parsed.Scheme == "" && parsed.Host == "" && strings.HasPrefix(parsed.Path, "/media/") && !existing {
			var asset models.ProductMedia
			if err := facades.Orm().Query().Where("url = ?", raw).First(&asset); err != nil || asset.ID == "" || actor == nil || !actor.Can("media.create", asset.OwnerID) {
				return errors.New("image asset not found or access denied")
			}
		}
		if !p.AllowRemote && parsed.Scheme == "https" && !existing {
			return errors.New("remote image links are disabled; upload an image instead")
		}
	}
	return nil
}
