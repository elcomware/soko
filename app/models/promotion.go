package models

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"
)

type CategorySelection struct {
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
}
type PromotionScope struct {
	All                  bool                `json:"all"`
	ProductIDs           []string            `json:"product_ids"`
	Categories           []string            `json:"categories"`
	SubCategories        []CategorySelection `json:"sub_categories"`
	ExcludeProductIDs    []string            `json:"exclude_product_ids"`
	ExcludeCategories    []string            `json:"exclude_categories"`
	ExcludeSubCategories []CategorySelection `json:"exclude_sub_categories"`
}
type Promotion struct {
	ID        string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Name      string         `gorm:"type:varchar(120)" json:"name"`
	VendorID  string         `gorm:"type:varchar(64);index" json:"vendor_id"` // empty means a platform campaign
	Kind      string         `gorm:"type:varchar(16)" json:"kind"`
	Amount    float64        `json:"amount"` // percentage points or USD per unit
	Enabled   bool           `json:"enabled"`
	StartsAt  time.Time      `json:"starts_at"`
	EndsAt    time.Time      `json:"ends_at"`
	ScopeJSON string         `gorm:"type:text" json:"-"`
	Scope     PromotionScope `gorm:"-" json:"scope"`
	CreatedBy string         `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Promotion) TableName() string { return "promotions" }
func (p *Promotion) Hydrate() {
	if p.ScopeJSON != "" {
		_ = json.Unmarshal([]byte(p.ScopeJSON), &p.Scope)
	}
}
func (p *Promotion) PrepareSave() { b, _ := json.Marshal(p.Scope); p.ScopeJSON = string(b) }
func (p Promotion) Validate() error {
	if strings.TrimSpace(p.Name) == "" || len(p.Name) > 120 || strings.ContainsAny(p.Name, "<>") {
		return errors.New("campaign name must be plain text between 1 and 120 characters")
	}
	if p.StartsAt.IsZero() || p.EndsAt.IsZero() || !p.EndsAt.After(p.StartsAt) {
		return errors.New("campaign end must be after its start; specify both times with a time zone")
	}
	if math.IsNaN(p.Amount) || math.IsInf(p.Amount, 0) || p.Amount <= 0 {
		return errors.New("discount must be a positive amount")
	}
	if (p.Kind != "percentage" && p.Kind != "fixed") || (p.Kind == "percentage" && p.Amount > 100) || (p.Kind == "fixed" && p.Amount > 1e12) {
		return errors.New("choose percentage (up to 100) or a fixed USD discount")
	}
	if math.Abs(p.Amount*100-math.Round(p.Amount*100)) > 0.00001 {
		return errors.New("discount supports at most two decimal places")
	}
	s := p.Scope
	if !s.All && len(s.ProductIDs)+len(s.Categories)+len(s.SubCategories) == 0 {
		return errors.New("select all products, a category, subcategory or individual product")
	}
	for _, values := range [][]string{s.ProductIDs, s.Categories, s.ExcludeProductIDs, s.ExcludeCategories} {
		if len(values) > 1000 {
			return errors.New("campaign target limit exceeded")
		}
		for _, v := range values {
			if strings.TrimSpace(v) == "" || len(v) > 128 {
				return errors.New("invalid campaign target")
			}
		}
	}
	for _, values := range [][]CategorySelection{s.SubCategories, s.ExcludeSubCategories} {
		if len(values) > 1000 {
			return errors.New("campaign target limit exceeded")
		}
		for _, v := range values {
			if strings.TrimSpace(v.Category) == "" || strings.TrimSpace(v.SubCategory) == "" || len(v.Category) > 128 || len(v.SubCategory) > 128 {
				return errors.New("subcategory targets need both category and subcategory")
			}
		}
	}
	return nil
}
func hasTarget(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
func hasSubcategory(values []CategorySelection, p Product) bool {
	for _, v := range values {
		if v.Category == p.Category && v.SubCategory == p.SubCategory {
			return true
		}
	}
	return false
}
func (c Promotion) Applies(p Product, now time.Time) bool {
	if !c.Enabled || now.Before(c.StartsAt) || !now.Before(c.EndsAt) || (c.VendorID != "" && c.VendorID != p.VendorID) {
		return false
	}
	s := c.Scope
	if hasTarget(s.ExcludeProductIDs, p.ID) || hasTarget(s.ExcludeCategories, p.Category) || hasSubcategory(s.ExcludeSubCategories, p) {
		return false
	}
	return s.All || hasTarget(s.ProductIDs, p.ID) || hasTarget(s.Categories, p.Category) || hasSubcategory(s.SubCategories, p)
}

type ProductPrice struct {
	Quantity          int        `json:"quantity"`
	Currency          string     `json:"currency"`
	OriginalUnitPrice float64    `json:"original_unit_price"`
	UnitPrice         float64    `json:"unit_price"`
	DiscountPerUnit   float64    `json:"discount_per_unit"`
	PromotionID       string     `json:"promotion_id,omitempty"`
	PromotionName     string     `json:"promotion_name,omitempty"`
	PromotionEndsAt   *time.Time `json:"promotion_ends_at,omitempty"`
}

// PriceProduct applies one best eligible campaign to the applicable quantity tier.
// Ties use campaign ID for deterministic results. All arithmetic after base-price
// validation uses cents, with half-up rounding; campaign discounts never stack.
func PriceProduct(p *Product, quantity int, campaigns []Promotion, now time.Time) (ProductPrice, error) {
	base, err := p.PriceFor(quantity)
	if err != nil {
		return ProductPrice{}, err
	}
	cents := int64(math.Round(base * 100))
	best := cents
	result := ProductPrice{Quantity: quantity, Currency: "USD", OriginalUnitPrice: base, UnitPrice: base}
	for _, c := range campaigns {
		c.Hydrate()
		if c.Validate() != nil || !c.Applies(*p, now) {
			continue
		}
		discount := int64(math.Round(c.Amount * 100))
		if c.Kind == "percentage" {
			basisPoints := int64(math.Round(c.Amount * 100))
			discount = (cents*basisPoints + 5000) / 10000
		}
		price := max(int64(0), cents-discount)
		if price < best || (price == best && price < cents && c.ID < result.PromotionID) {
			best = price
			result.PromotionID = c.ID
			result.PromotionName = c.Name
			end := c.EndsAt
			result.PromotionEndsAt = &end
		}
	}
	result.UnitPrice = float64(best) / 100
	result.DiscountPerUnit = float64(cents-best) / 100
	return result, nil
}
