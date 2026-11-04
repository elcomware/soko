package actions

import (
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"errors"
	"github.com/google/uuid"
	"math"
	"net/url"
	"strings"
	"time"
)

type CreateProductAction struct{}

func NewCreateProductAction() *CreateProductAction { return &CreateProductAction{} }
func ValidateProduct(r *requests.CreateProductRequest) error {
	r.Region = models.NormalizeRegion(r.Region)
	if len(strings.TrimSpace(r.Title)) < 3 || len(r.Title) > 255 || strings.TrimSpace(r.Category) == "" || len(r.Category) > 128 {
		return errors.New("provide a title and category within their length limits")
	}
	if math.IsNaN(r.Price) || math.IsInf(r.Price, 0) || math.Round(r.Price*100) <= 0 || r.Price > 1e12 || math.IsNaN(r.WholesalePrice) || math.IsInf(r.WholesalePrice, 0) || r.WholesalePrice < 0 || r.WholesalePrice > r.Price {
		return errors.New("provide a positive price and a wholesale price no higher than it")
	}
	if r.MOQ < 1 || r.MOQ > 1000000 || r.Stock < -1 || r.Stock > 1000000000 {
		return errors.New("minimum quantity must be positive and stock within its limits")
	}
	if r.Type != "physical" && r.Type != "digital_software" && r.Type != "digital_asset" && r.Type != "digital_license" {
		return errors.New("invalid product type")
	}
	if r.Format != "retail" && r.Format != "wholesale" && r.Format != "digital" && r.Format != "make_offer" {
		return errors.New("invalid listing format")
	}
	if r.Condition != "new" && r.Condition != "refurbished" && r.Condition != "open_box" {
		return errors.New("invalid item condition")
	}
	if r.Type == "physical" && (r.Stock < 0 || r.Format == "digital") {
		return errors.New("physical listings need nonnegative stock and a physical format")
	}
	if len(r.Description) > 20000 || len(r.Sku) > 64 || len(r.SubCategory) > 128 || len(r.DigitalFileSize) > 64 || len(r.DigitalLicenseFormat) > 64 {
		return errors.New("product details exceed their length limits")
	}
	validRegion := map[string]bool{"": true, "global": true, "east-africa": true, "west-africa": true, "north-america": true, "europe": true, "asia-pacific": true, "middle-east": true}
	if !validRegion[r.Region] {
		return errors.New("invalid market region")
	}
	for _, raw := range []string{r.ImageUrl, r.DigitalDownloadUrl} {
		if raw != "" {
			u, err := url.Parse(raw)
			if len(raw) > 2048 || err != nil || !((u.Scheme == "https" && u.Host != "") || (u.Scheme == "" && strings.HasPrefix(raw, "/") && !strings.HasPrefix(raw, "//"))) {
				return errors.New("URLs must use HTTPS or a local absolute path")
			}
		}
	}
	if len(r.PriceTiers) > 20 {
		return errors.New("volume pricing supports at most 20 tiers")
	}
	previousMax := 0
	for i, tier := range r.PriceTiers {
		if tier.MinQty < r.MOQ || tier.MinQty > 1000000 || tier.MaxQty < 0 || tier.MaxQty > 1000000 || (tier.MaxQty != 0 && tier.MaxQty < tier.MinQty) || (i > 0 && tier.MinQty <= previousMax) || (tier.MaxQty == 0 && i < len(r.PriceTiers)-1) {
			return errors.New("volume tiers must be ordered, nonoverlapping quantity ranges at or above the listing minimum")
		}
		if math.IsNaN(tier.Price) || math.IsInf(tier.Price, 0) || tier.Price > 1e12 || math.Round(tier.Price*100) < 1 {
			return errors.New("each volume tier needs a positive finite price of at least one cent")
		}
		r.PriceTiers[i].Price = math.Round(tier.Price*100) / 100
		previousMax = tier.MaxQty
	}
	return nil
}
func ApplyProductInput(p *models.Product, r *requests.CreateProductRequest) {
	p.Title = strings.TrimSpace(r.Title)
	p.Category = strings.TrimSpace(r.Category)
	p.SubCategory = r.SubCategory
	p.Description = r.Description
	p.ShortDesc = r.Description
	p.Price = r.Price
	p.WholesalePrice = r.WholesalePrice
	p.MOQ = r.MOQ
	p.Stock = r.Stock
	p.Sku = strings.TrimSpace(r.Sku)
	p.ImageUrl = r.ImageUrl
	p.Type = models.ProductType(r.Type)
	p.Format = models.ListingFormat(r.Format)
	p.Condition = models.ItemCondition(r.Condition)
	p.VendorRegion = models.NormalizeRegion(r.Region)
	p.AcceptsOffers = r.AcceptsOffers
	p.DigitalDownloadUrl = r.DigitalDownloadUrl
	p.DigitalLicenseFormat = r.DigitalLicenseFormat
	p.DigitalFileSize = r.DigitalFileSize
	p.UpdatedAt = time.Now()
	if p.VendorRegion == "" {
		p.VendorRegion = "global"
	}
	if p.IsDigital() {
		p.Format = models.FormatDigital
		p.LeadTimeDays = 0
		p.Incoterm = "Digital delivery after payment confirmation"
	} else {
		p.DigitalDownloadUrl = ""
		p.DigitalLicenseFormat = ""
		p.DigitalFileSize = ""
	}
	if r.PriceTiers != nil {
		p.PriceTiers = append([]models.PriceTier{}, r.PriceTiers...)
	}
	p.PrepareSave()
}
func (a *CreateProductAction) ExecuteFromRequest(r *requests.CreateProductRequest, u *models.User) (*models.Product, error) {
	if !u.CanManageCatalog() {
		return nil, errors.New("catalog access is restricted to suppliers and administrators")
	}
	if err := ValidateProduct(r); err != nil {
		return nil, err
	}
	p := &models.Product{ID: "prod_" + uuid.NewString(), VendorID: u.ID, VendorName: u.CompanyName, VendorCountry: u.Country, VendorVerified: u.Verified, VendorBadge: "Supplier", CreatedAt: time.Now()}
	if p.VendorName == "" {
		p.VendorName = u.Name
	}
	ApplyProductInput(p, r)
	p.Slug = strings.ToLower(strings.ReplaceAll(p.Title, " ", "-")) + "-" + p.ID
	if err := facades.Orm().Query().Create(p); err != nil {
		return nil, err
	}
	return p, nil
}
