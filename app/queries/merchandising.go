package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func DefaultSectionOrder() []string {
	return []string{"showcase", "category_navigation", "catalog", "discovery", "sourcing"}
}
func ValidateSectionOrder(order []string) error {
	expected := DefaultSectionOrder()
	if len(order) != len(expected) {
		return errors.New("include each storefront section exactly once; the product catalog is required")
	}
	allowed := map[string]bool{}
	for _, key := range expected {
		allowed[key] = true
	}
	for _, key := range order {
		if !allowed[key] {
			return errors.New("unknown or duplicate storefront section")
		}
		delete(allowed, key)
	}
	return nil
}
func DefaultMerchandising() models.MerchandisingSettings {
	return models.MerchandisingSettings{SectionOrder: DefaultSectionOrder(), Announcement: "A world of suppliers. One marketplace.", Hero: models.MerchandisingPlacement{Enabled: true, Title: "Selected for your next order", Strategy: "featured", ProductIDs: []string{}, Fallback: "newest", Limit: 3}, HeroSide: models.MerchandisingPlacement{Enabled: true, Title: "More to discover", Strategy: "newest", ProductIDs: []string{}, Fallback: "newest", Limit: 2}, Discovery: models.MerchandisingPlacement{Enabled: true, Title: "Digital tools for your next project", Strategy: "newest", ProductIDs: []string{}, Fallback: "newest", Limit: 4}, CategoryNavigation: models.CategoryNavigation{Enabled: true, Title: "Explore categories"}}
}
func LoadMerchandising() (models.MerchandisingSettings, error) {
	defaults := DefaultMerchandising()
	var record models.MerchandisingRecord
	if err := facades.Orm().Query().Where("id = ?", "store").First(&record); err != nil {
		return defaults, err
	}
	if record.ID == "" {
		defaults.Announcement = LoadStorefrontContent()["announcement"]
		return defaults, nil
	}
	result := defaults
	if err := json.Unmarshal([]byte(record.SettingsJSON), &result); err != nil {
		return defaults, err
	}
	if err := ValidateMerchandising(result); err != nil {
		return defaults, err
	}
	result.Version = record.Version
	return result, nil
}
func plainMerchandisingText(value string, max int) bool {
	if !utf8.ValidString(value) || utf8.RuneCountInString(value) > max || strings.ContainsAny(value, "<>") {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func ValidateMerchandising(c models.MerchandisingSettings) error {
	if err := ValidateSectionOrder(c.SectionOrder); err != nil {
		return err
	}
	if !plainMerchandisingText(c.Announcement, 160) || c.Version < 0 || !plainMerchandisingText(c.CategoryNavigation.Title, 100) {
		return errors.New("invalid merchandising version or navigation title")
	}
	for i, p := range []models.MerchandisingPlacement{c.Hero, c.HeroSide, c.Discovery} {
		if !plainMerchandisingText(p.Title, 100) || !plainMerchandisingText(p.Category, 128) {
			return errors.New("titles and categories must be bounded plain text")
		}
		max := 12
		if i == 1 {
			max = 2
		}
		if p.Limit < 1 || p.Limit > max {
			return errors.New("invalid placement product limit")
		}
		if p.Strategy != "manual" && p.Strategy != "featured" && p.Strategy != "newest" && p.Strategy != "promotions" {
			return errors.New("unsupported product selection strategy")
		}
		if p.Fallback != "hide" && p.Fallback != "newest" {
			return errors.New("choose hide or newest fallback")
		}
		if len(p.ProductIDs) > 24 {
			return errors.New("at most 24 product pins per placement")
		}
		if !p.Enabled || p.Strategy != "manual" {
			continue
		}
		seen := map[string]bool{}
		for _, id := range p.ProductIDs {
			if id == "" || len(id) > 64 || strings.TrimSpace(id) != id || seen[id] {
				return errors.New("invalid or duplicate pinned product")
			}
			seen[id] = true
		}
	}
	return nil
}

type ResolvedPlacement struct {
	Strategy        string           `json:"strategy"`
	Enabled         bool             `json:"enabled"`
	Title           string           `json:"title"`
	Products        []models.Product `json:"products"`
	RotationSeconds int              `json:"rotation_seconds,omitempty"`
}
type ResolvedMerchandising struct {
	SectionOrder       []string                  `json:"section_order"`
	Hero               ResolvedPlacement         `json:"hero"`
	HeroSide           ResolvedPlacement         `json:"hero_side"`
	Discovery          ResolvedPlacement         `json:"discovery"`
	CategoryNavigation models.CategoryNavigation `json:"category_navigation"`
}

// ResolveMerchandising uses ONLY the already filtered public catalog supplied by
// the caller; pins never escape a shopper's current region/category/search.
func ResolveMerchandising(c models.MerchandisingSettings, products []models.Product, now time.Time) ResolvedMerchandising {
	if ValidateMerchandising(c) != nil {
		c = DefaultMerchandising()
	}
	eligible := []models.Product{}
	for _, p := range products {
		quantity := p.MOQ
		if quantity < 1 {
			quantity = 1
		}
		price, err := p.PriceFor(quantity)
		if p.ID != "" && strings.TrimSpace(p.Title) != "" && err == nil && !math.IsNaN(price) && !math.IsInf(price, 0) && price >= 0 && (p.Stock < 0 || p.Stock >= quantity) {
			eligible = append(eligible, p)
		}
	}
	sort.Slice(eligible, func(i, j int) bool {
		if eligible[i].CreatedAt.Equal(eligible[j].CreatedAt) {
			return eligible[i].ID < eligible[j].ID
		}
		return eligible[i].CreatedAt.After(eligible[j].CreatedAt)
	})
	selectProducts := func(config models.MerchandisingPlacement, digital bool, exclude map[string]bool) ResolvedPlacement {
		result := ResolvedPlacement{Enabled: config.Enabled, Title: config.Title, Strategy: config.Strategy, Products: []models.Product{}}
		if !config.Enabled {
			return result
		}
		pool := []models.Product{}
		for _, p := range eligible {
			if exclude[p.ID] || (digital && !p.IsDigital()) || (config.Category != "" && p.Category != config.Category) {
				continue
			}
			pool = append(pool, p)
		}
		chosen := map[string]bool{}
		add := func(p models.Product) {
			if !chosen[p.ID] {
				result.Products = append(result.Products, p.Public())
				chosen[p.ID] = true
			}
		}
		if config.Strategy == "manual" {
			for _, id := range config.ProductIDs {
				for _, p := range pool {
					if p.ID == id {
						add(p)
					}
				}
			}
		} else {
			for _, p := range pool {
				match := config.Strategy == "newest" || (config.Strategy == "featured" && p.Featured)
				if config.Strategy == "promotions" {
					match = p.Pricing != nil && p.Pricing.PromotionID != "" && p.Pricing.DiscountPerUnit > 0 && (p.Pricing.PromotionEndsAt == nil || p.Pricing.PromotionEndsAt.After(now))
				}
				if match {
					add(p)
				}
			}
		}
		if config.Fallback == "newest" {
			for _, p := range pool {
				add(p)
			}
		}
		if config.Strategy != "manual" {
			// Rank before truncation, including fallback candidates. A missing or
			// quarantined image must not crowd usable product media out of a slot.
			usable := func(p models.Product) bool {
				return strings.TrimSpace(p.ImageUrl) != "" && !strings.Contains(p.ImageUrl, "photo-1500937386664-56d1dfef3854")
			}
			sort.SliceStable(result.Products, func(i, j int) bool { return usable(result.Products[i]) && !usable(result.Products[j]) })
		}
		if len(result.Products) > config.Limit {
			result.Products = result.Products[:config.Limit]
		}
		if len(result.Products) == 0 {
			result.Enabled = false
		}
		return result
	}
	hero := selectProducts(c.Hero, false, nil)
	used := map[string]bool{}
	for _, p := range hero.Products {
		used[p.ID] = true
	}
	return ResolvedMerchandising{SectionOrder: append([]string(nil), c.SectionOrder...), Hero: hero, HeroSide: selectProducts(c.HeroSide, false, used), Discovery: selectProducts(c.Discovery, true, nil), CategoryNavigation: c.CategoryNavigation}
}
