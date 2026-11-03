package models

import "time"

type ProductType string

const (
	TypePhysical        ProductType = "physical"
	TypeDigitalSoftware ProductType = "digital_software"
	TypeDigitalAsset    ProductType = "digital_asset"
	TypeDigitalLicense  ProductType = "digital_license"
)

type ListingFormat string

const (
	FormatWholesale ListingFormat = "wholesale"  // Alibaba style (MOQ + Volume Tiers)
	FormatRetail    ListingFormat = "retail"     // Amazon style (Buy Now + Fast Dispatch)
	FormatDigital   ListingFormat = "digital"    // Instant Automated License / Download
	FormatMakeOffer ListingFormat = "make_offer" // eBay style (Accepts Offers / Negotiation)
)

type ItemCondition string

const (
	ConditionNew         ItemCondition = "new"
	ConditionRefurbished ItemCondition = "refurbished"
	ConditionOpenBox     ItemCondition = "open_box"
)

type PriceTier struct {
	MinQty int     `json:"min_qty"`
	MaxQty int     `json:"max_qty"`
	Price  float64 `json:"price"`
}

type Product struct {
	ID                   string        `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Title                string        `gorm:"type:varchar(255)" json:"title"`
	Slug                 string        `gorm:"type:varchar(255)" json:"slug"`
	Description          string        `gorm:"type:text" json:"description"`
	ShortDesc            string        `gorm:"type:text" json:"short_desc"`
	Category             string        `gorm:"type:varchar(128);index" json:"category"`
	SubCategory          string        `gorm:"type:varchar(128)" json:"sub_category"`
	Type                 ProductType   `gorm:"type:varchar(64);index" json:"type"`
	Format               ListingFormat `gorm:"type:varchar(64);index" json:"format"`
	Condition            ItemCondition `gorm:"type:varchar(64)" json:"condition"`
	Price                float64       `gorm:"type:double" json:"price"`
	WholesalePrice       float64       `gorm:"type:double" json:"wholesale_price"`
	PriceTiersJSON       string        `gorm:"column:price_tiers;type:text" json:"-"`
	MOQ                  int           `gorm:"column:moq;type:integer" json:"moq"`
	Stock                int           `gorm:"type:integer" json:"stock"`
	Sku                  string        `gorm:"type:varchar(64)" json:"sku"`
	VendorID             string        `gorm:"type:varchar(64)" json:"vendor_id"`
	VendorName           string        `gorm:"type:varchar(128)" json:"vendor_name"`
	VendorCountry        string        `gorm:"type:varchar(64)" json:"vendor_country"`
	VendorRegion         string        `gorm:"type:varchar(64);index" json:"vendor_region"`
	VendorVerified       bool          `gorm:"type:boolean" json:"vendor_verified"`
	VendorBadge          string        `gorm:"type:varchar(64)" json:"vendor_badge"`
	VendorPositiveRate   string        `gorm:"type:varchar(32)" json:"vendor_positive_rate"`
	TradeAssuranceCover  float64       `gorm:"type:double" json:"trade_assurance_cover"`
	AcceptsOffers        bool          `gorm:"type:boolean" json:"accepts_offers"`
	Rating               float64       `gorm:"type:double" json:"rating"`
	ReviewsCount         int           `gorm:"type:integer" json:"reviews_count"`
	ImageUrl             string        `gorm:"type:text" json:"image_url"`
	GalleryUrlsJSON      string        `gorm:"column:gallery_urls;type:text" json:"-"`
	TagsJSON             string        `gorm:"column:tags;type:text" json:"-"`
	Featured             bool          `gorm:"type:boolean;index" json:"featured"`
	LeadTimeDays         int           `gorm:"type:integer" json:"lead_time_days"`
	Incoterm             string        `gorm:"type:varchar(32)" json:"incoterm"`
	LocalPickupAvailable bool          `gorm:"type:boolean" json:"local_pickup_available"`
	DigitalDownloadUrl   string        `gorm:"type:text" json:"digital_download_url,omitempty"`
	DigitalLicenseFormat string        `gorm:"type:varchar(64)" json:"digital_license_format,omitempty"`
	DigitalFileSize      string        `gorm:"type:varchar(64)" json:"digital_file_size,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`

	// In-memory or mapped fields
	PriceTiers           []PriceTier   `gorm:"-" json:"price_tiers,omitempty"`
	GalleryUrls          []string      `gorm:"-" json:"gallery_urls,omitempty"`
	Tags                 []string      `gorm:"-" json:"tags,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
