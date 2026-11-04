package queries

import (
	"errors"
	"strings"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type ProductFilter struct {
	Type      string
	Format    string
	Region    string
	Category  string
	Condition string
	Query     string
}

type ProductQueries struct{}

func NewProductQueries() *ProductQueries {
	return &ProductQueries{}
}

func (q *ProductQueries) FindFiltered(filter ProductFilter) ([]models.Product, error) {
	builder := facades.Orm().Query().Model(&models.Product{})
	filter.Region = models.NormalizeRegion(filter.Region)

	if filter.Type != "" && filter.Type != "all" {
		builder = builder.Where("type = ?", filter.Type)
	}
	if filter.Format != "" && filter.Format != "all" {
		builder = builder.Where("format = ?", filter.Format)
	}
	if filter.Region != "" && filter.Region != "global" && filter.Region != "all" {
		aliases := models.RegionAliases(filter.Region)
		placeholders := make([]string, len(aliases))
		args := make([]any, len(aliases))
		for i, alias := range aliases {
			placeholders[i] = "?"
			args[i] = alias
		}
		builder = builder.Where("LOWER(TRIM(vendor_region)) IN ("+strings.Join(placeholders, ",")+")", args...)
	}
	if filter.Category != "" && filter.Category != "all" {
		builder = builder.Where("category = ?", filter.Category)
	}
	if filter.Condition != "" && filter.Condition != "all" {
		builder = builder.Where("condition = ?", filter.Condition)
	}
	if filter.Query != "" {
		search := "%" + strings.ToLower(filter.Query) + "%"
		builder = builder.Where("(LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(category) LIKE ? OR LOWER(sku) LIKE ?)",
			search, search, search, search)
	}

	var products []models.Product
	if err := builder.OrderByDesc("created_at").Get(&products); err != nil {
		return nil, err
	}

	for i := range products {
		products[i].Hydrate()
	}

	return products, nil
}

func (q *ProductQueries) FindByID(id string) (*models.Product, error) {
	var product models.Product
	if err := facades.Orm().Query().Where("id = ?", id).First(&product); err != nil {
		return nil, err
	}
	if product.ID == "" {
		return nil, errors.New("listing not found")
	}
	product.Hydrate()
	return &product, nil
}

func (q *ProductQueries) Count() (int64, error) {
	return facades.Orm().Query().Model(&models.Product{}).Count()
}
