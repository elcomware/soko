package models

import "time"

type ProductType string

const (
	TypePhysical        ProductType = "physical"
	TypeDigitalSoftware ProductType = "digital_software"
	TypeDigitalAsset    ProductType = "digital_asset"
	TypeDigitalLicense  ProductType = "digital_license"
)

type Product struct {
	ID                   string      `json:"id"`
	Title                string      `json:"title"`
	Slug                 string      `json:"slug"`
	Description          string      `json:"description"`
	ShortDesc            string      `json:"short_desc"`
	Category             string      `json:"category"`
	SubCategory          string      `json:"sub_category"`
	Type                 ProductType `json:"type"` // physical, digital_software, digital_asset, digital_license
	Price                float64     `json:"price"` // Retail unit price
	WholesalePrice       float64     `json:"wholesale_price"` // B2B tier unit price
	MOQ                  int         `json:"moq"` // Minimum Order Quantity (1 for digital)
	Stock                int         `json:"stock"` // Available stock (-1 for unlimited digital)
	Sku                  string      `json:"sku"`
	VendorID             string      `json:"vendor_id"`
	VendorName           string      `json:"vendor_name"`
	VendorCountry        string      `json:"vendor_country"`
	VendorVerified       bool        `json:"vendor_verified"`
	Rating               float64     `json:"rating"`
	ReviewsCount         int         `json:"reviews_count"`
	ImageUrl             string      `json:"image_url"`
	Tags                 []string    `json:"tags"`
	Featured             bool        `json:"featured"`
	LeadTimeDays         int         `json:"lead_time_days"` // 0 for instant digital, >0 for freight
	DigitalDownloadUrl   string      `json:"digital_download_url,omitempty"`
	DigitalLicenseFormat string      `json:"digital_license_format,omitempty"`
	DigitalFileSize      string      `json:"digital_file_size,omitempty"`
	CreatedAt            time.Time   `json:"created_at"`
}

