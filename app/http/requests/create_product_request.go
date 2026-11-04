package requests

import (
	"elcomware/app/models"
	"github.com/goravel/framework/contracts/http"
)

type CreateProductRequest struct {
	PriceTiers           []models.PriceTier `form:"price_tiers" json:"price_tiers"`
	Version              string             `form:"version" json:"version"`
	Title                string             `form:"title" json:"title"`
	Category             string             `form:"category" json:"category"`
	SubCategory          string             `form:"sub_category" json:"sub_category"`
	Type                 string             `form:"type" json:"type"`
	Format               string             `form:"format" json:"format"`
	Condition            string             `form:"condition" json:"condition"`
	Region               string             `form:"region" json:"region"`
	Price                float64            `form:"price" json:"price"`
	WholesalePrice       float64            `form:"wholesale_price" json:"wholesale_price"`
	MOQ                  int                `form:"moq" json:"moq"`
	Stock                int                `form:"stock" json:"stock"`
	Sku                  string             `form:"sku" json:"sku"`
	Description          string             `form:"description" json:"description"`
	ImageUrl             string             `form:"image_url" json:"image_url"`
	DigitalDownloadUrl   string             `form:"digital_download_url" json:"digital_download_url"`
	DigitalLicenseFormat string             `form:"digital_license_format" json:"digital_license_format"`
	DigitalFileSize      string             `form:"digital_file_size" json:"digital_file_size"`
	AcceptsOffers        bool               `form:"accepts_offers" json:"accepts_offers"`
}

func (r *CreateProductRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CreateProductRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"title":    "required|min_len:3",
		"category": "required",
		"price":    "required",
	}
}

func (r *CreateProductRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"title.required":    "Product title is required.",
		"title.min_len":     "Product title must have at least 3 characters.",
		"category.required": "Marketplace category is required.",
		"price.required":    "Unit base price is required.",
	}
}
