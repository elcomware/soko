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
	ID                   string        `json:"id"`
	Title                string        `json:"title"`
	Slug                 string        `json:"slug"`
	Description          string        `json:"description"`
	ShortDesc            string        `json:"short_desc"`
	Category             string        `json:"category"`
	SubCategory          string        `json:"sub_category"`
	Type                 ProductType   `json:"type"`                  // physical, digital_software, digital_asset, digital_license
	Format               ListingFormat `json:"format"`                // wholesale, retail, digital, make_offer
	Condition            ItemCondition `json:"condition"`             // new, refurbished, open_box
	Price                float64       `json:"price"`                 // Retail / base unit price
	WholesalePrice       float64       `json:"wholesale_price"`       // Default B2B wholesale unit price
	PriceTiers           []PriceTier   `json:"price_tiers,omitempty"` // Alibaba-style tiered volume brackets
	MOQ                  int           `json:"moq"`                   // Minimum Order Quantity
	Stock                int           `json:"stock"`                 // Available inventory (-1 for unlimited digital)
	Sku                  string        `json:"sku"`
	VendorID             string        `json:"vendor_id"`
	VendorName           string        `json:"vendor_name"`
	VendorCountry        string        `json:"vendor_country"`
	VendorRegion         string        `json:"vendor_region"` // North America, Europe, East Africa, West Africa, Asia-Pacific
	VendorVerified       bool          `json:"vendor_verified"`
	VendorBadge          string        `json:"vendor_badge"`          // Gold Verified, Top Rated Plus, Certified OEM
	VendorPositiveRate   string        `json:"vendor_positive_rate"`  // e.g. "99.6%"
	TradeAssuranceCover  float64       `json:"trade_assurance_cover"` // e.g. $1,000,000 escrow protection
	AcceptsOffers        bool          `json:"accepts_offers"`        // eBay-style Make an Offer enabled
	Rating               float64       `json:"rating"`
	ReviewsCount         int           `json:"reviews_count"`
	ImageUrl             string        `json:"image_url"`
	GalleryUrls          []string      `json:"gallery_urls,omitempty"`
	Tags                 []string      `json:"tags"`
	Featured             bool          `json:"featured"`
	LeadTimeDays         int           `json:"lead_time_days"` // 0 for instant digital, >0 for freight
	Incoterm             string        `json:"incoterm"`       // DDP, FOB, CIF, Instant Delivery
	LocalPickupAvailable bool          `json:"local_pickup_available"`
	DigitalDownloadUrl   string        `json:"digital_download_url,omitempty"`
	DigitalLicenseFormat string        `json:"digital_license_format,omitempty"`
	DigitalFileSize      string        `json:"digital_file_size,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
}
