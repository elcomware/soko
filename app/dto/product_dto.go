package dto

import (
	"time"

	"elcomware/app/models"
)

type ProductFilterInput struct {
	Type      string `json:"type"`
	Format    string `json:"format"`
	Region    string `json:"region"`
	Category  string `json:"category"`
	Condition string `json:"condition"`
	Query     string `json:"query"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}

type CreateProductInput struct {
	Title                string             `json:"title"`
	Slug                 string             `json:"slug"`
	Description          string             `json:"description"`
	ShortDesc            string             `json:"short_desc"`
	Category             string             `json:"category"`
	SubCategory          string             `json:"sub_category"`
	Type                 models.ProductType `json:"type"`
	Format               models.ListingFormat `json:"format"`
	Condition            models.ItemCondition `json:"condition"`
	Price                float64            `json:"price"`
	WholesalePrice       float64            `json:"wholesale_price"`
	PriceTiers           []models.PriceTier `json:"price_tiers"`
	MOQ                  int                `json:"moq"`
	Stock                int                `json:"stock"`
	Sku                  string             `json:"sku"`
	VendorID             string             `json:"vendor_id"`
	VendorName           string             `json:"vendor_name"`
	VendorCountry        string             `json:"vendor_country"`
	VendorRegion         string             `json:"vendor_region"`
	VendorVerified       bool               `json:"vendor_verified"`
	VendorBadge          string             `json:"vendor_badge"`
	VendorPositiveRate   string             `json:"vendor_positive_rate"`
	TradeAssuranceCover  float64            `json:"trade_assurance_cover"`
	AcceptsOffers        bool               `json:"accepts_offers"`
	ImageUrl             string             `json:"image_url"`
	GalleryUrls          []string           `json:"gallery_urls"`
	Tags                 []string           `json:"tags"`
	Featured             bool               `json:"featured"`
	LeadTimeDays         int                `json:"lead_time_days"`
	Incoterm             string             `json:"incoterm"`
	LocalPickupAvailable bool               `json:"local_pickup_available"`
	DigitalDownloadUrl   string             `json:"digital_download_url"`
	DigitalLicenseFormat string             `json:"digital_license_format"`
	DigitalFileSize      string             `json:"digital_file_size"`
}

type ProductOutput struct {
	ID                   string               `json:"id"`
	Title                string               `json:"title"`
	Slug                 string               `json:"slug"`
	Description          string               `json:"description"`
	ShortDesc            string               `json:"short_desc"`
	Category             string               `json:"category"`
	SubCategory          string               `json:"sub_category"`
	Type                 models.ProductType   `json:"type"`
	Format               models.ListingFormat `json:"format"`
	Condition            models.ItemCondition `json:"condition"`
	Price                float64              `json:"price"`
	WholesalePrice       float64              `json:"wholesale_price"`
	PriceTiers           []models.PriceTier   `json:"price_tiers,omitempty"`
	MOQ                  int                  `json:"moq"`
	Stock                int                  `json:"stock"`
	Sku                  string               `json:"sku"`
	VendorID             string               `json:"vendor_id"`
	VendorName           string               `json:"vendor_name"`
	VendorCountry        string               `json:"vendor_country"`
	VendorRegion         string               `json:"vendor_region"`
	VendorVerified       bool                 `json:"vendor_verified"`
	VendorBadge          string               `json:"vendor_badge"`
	VendorPositiveRate   string               `json:"vendor_positive_rate"`
	TradeAssuranceCover  float64              `json:"trade_assurance_cover"`
	AcceptsOffers        bool                 `json:"accepts_offers"`
	Rating               float64              `json:"rating"`
	ReviewsCount         int                  `json:"reviews_count"`
	ImageUrl             string               `json:"image_url"`
	GalleryUrls          []string             `json:"gallery_urls,omitempty"`
	Tags                 []string             `json:"tags"`
	Featured             bool                 `json:"featured"`
	LeadTimeDays         int                  `json:"lead_time_days"`
	Incoterm             string               `json:"incoterm"`
	LocalPickupAvailable bool                 `json:"local_pickup_available"`
	DigitalDownloadUrl   string               `json:"digital_download_url,omitempty"`
	DigitalLicenseFormat string               `json:"digital_license_format,omitempty"`
	DigitalFileSize      string               `json:"digital_file_size,omitempty"`
	CreatedAt            time.Time            `json:"created_at"`
}
