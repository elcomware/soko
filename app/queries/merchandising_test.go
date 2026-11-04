package queries

import (
	"elcomware/app/models"
	"encoding/json"
	"testing"
	"time"
)

func merchProduct(id string) models.Product {
	return models.Product{ID: id, Title: id, Price: 10, Stock: 10, MOQ: 1, Type: models.TypePhysical, CreatedAt: time.Unix(100, 0)}
}
func TestMerchandisingSelectionEligibilityAndNoHeroDuplicates(t *testing.T) {
	c := DefaultMerchandising()
	c.Hero.Strategy = "manual"
	c.Hero.ProductIDs = []string{"outside-filter", "b", "a"}
	c.Hero.Limit = 2
	c.Hero.Fallback = "hide"
	products := []models.Product{merchProduct("a"), merchProduct("b"), merchProduct("c"), merchProduct("soldout")}
	products[3].Stock = 0
	result := ResolveMerchandising(c, products, time.Now())
	if len(result.Hero.Products) != 2 || result.Hero.Products[0].ID != "b" || result.Hero.Products[1].ID != "a" {
		t.Fatal("manual/filter selection failed", result.Hero)
	}
	if len(result.HeroSide.Products) != 1 || result.HeroSide.Products[0].ID != "c" {
		t.Fatal("duplicate or unavailable side card", result.HeroSide)
	}
	if result.Discovery.Enabled {
		t.Fatal("physical products in digital shelf")
	}
	empty := ResolveMerchandising(c, nil, time.Now())
	if empty.Hero.Enabled || empty.HeroSide.Enabled {
		t.Fatal("empty showcases remain enabled")
	}
}
func TestMerchandisingPromotionExpirationFallbackAndPublicData(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Second)
	future := now.Add(time.Hour)
	expired := merchProduct("expired")
	expired.Pricing = &models.ProductPrice{PromotionID: "old", DiscountPerUnit: 2, PromotionEndsAt: &past}
	active := merchProduct("active")
	active.Type = models.TypeDigitalSoftware
	active.Stock = -1
	active.DigitalDownloadUrl = "private"
	active.DigitalLicenseFormat = "secret"
	active.Pricing = &models.ProductPrice{PromotionID: "live", DiscountPerUnit: 1, PromotionEndsAt: &future}
	c := DefaultMerchandising()
	c.Hero.Strategy = "promotions"
	c.Hero.Fallback = "hide"
	result := ResolveMerchandising(c, []models.Product{expired, active}, now)
	if len(result.Hero.Products) != 1 || result.Hero.Products[0].ID != "active" || result.Hero.Products[0].DigitalDownloadUrl != "" || result.Hero.Products[0].DigitalLicenseFormat != "" {
		t.Fatal("promotion/privacy selection failed")
	}
	if len(result.Discovery.Products) != 1 {
		t.Fatal("digital listing omitted")
	}
	c.Hero.ProductIDs = []string{"missing"}
	c.Hero.Strategy = "manual"
	c.Hero.Fallback = "newest"
	if len(ResolveMerchandising(c, []models.Product{expired}, now).Hero.Products) != 1 {
		t.Fatal("fallback absent")
	}
}
func TestMerchandisingValidation(t *testing.T) {
	valid := DefaultMerchandising()
	if ValidateMerchandising(valid) != nil {
		t.Fatal("bad defaults")
	}
	for _, mutate := range []func(*models.MerchandisingSettings){func(c *models.MerchandisingSettings) { c.HeroSide.Limit = 3 }, func(c *models.MerchandisingSettings) { c.Hero.Strategy = "best_selling" }, func(c *models.MerchandisingSettings) { c.Hero.Title = "<script>" }, func(c *models.MerchandisingSettings) {
		c.Hero.Strategy = "manual"
		c.Hero.ProductIDs = []string{"a", "a"}
	}, func(c *models.MerchandisingSettings) { c.Hero.Fallback = "random" }} {
		c := valid
		mutate(&c)
		if ValidateMerchandising(c) == nil {
			t.Fatal("invalid settings accepted")
		}
	}
}

func TestAutomaticMerchandisingRanksImagesBeforeLimitButManualKeepsOrder(t *testing.T) {
	missing := merchProduct("a")
	wrong := merchProduct("b")
	wrong.ImageUrl = "https://images.unsplash.com/photo-1500937386664-56d1dfef3854"
	usable := merchProduct("z")
	usable.ImageUrl = "/storage/product-media/actual-product.png"
	c := DefaultMerchandising()
	c.Hero.Strategy = "newest"
	c.Hero.Limit = 1
	resolved := ResolveMerchandising(c, []models.Product{missing, wrong, usable}, time.Now())
	if resolved.Hero.Products[0].ID != "z" || resolved.Hero.Strategy != "newest" {
		t.Fatal("usable media excluded by pre-ranking limit", resolved.Hero)
	}
	c.Hero.Strategy = "manual"
	c.Hero.ProductIDs = []string{"b", "z"}
	resolved = ResolveMerchandising(c, []models.Product{missing, wrong, usable}, time.Now())
	if resolved.Hero.Products[0].ID != "b" || resolved.Hero.Strategy != "manual" {
		t.Fatal("manual pin order changed", resolved.Hero)
	}
}

func TestSectionOrderIsClosedCompleteAndPreservesPublishedSequence(t *testing.T) {
	for _, order := range [][]string{nil, {}, {"catalog"}, {"showcase", "category_navigation", "catalog", "discovery", "catalog"}, {"showcase", "category_navigation", "catalog", "discovery", "<script>"}} {
		if ValidateSectionOrder(order) == nil {
			t.Fatal("unsafe/incomplete order accepted", order)
		}
	}
	c := DefaultMerchandising()
	c.SectionOrder = []string{"catalog", "showcase", "discovery", "category_navigation", "sourcing"}
	if err := ValidateMerchandising(c); err != nil {
		t.Fatal(err)
	}
	result := ResolveMerchandising(c, nil, time.Now())
	if result.SectionOrder[0] != "catalog" {
		t.Fatal("custom sequence discarded")
	}
	result.SectionOrder[0] = "changed"
	if c.SectionOrder[0] != "catalog" {
		t.Fatal("resolved sequence aliases stored settings")
	}
}

func TestLegacyMerchandisingJSONInheritsDefaultSectionOrder(t *testing.T) {
	settings := DefaultMerchandising()
	if err := json.Unmarshal([]byte(`{"announcement":"Existing merchant announcement"}`), &settings); err != nil {
		t.Fatal(err)
	}
	if err := ValidateMerchandising(settings); err != nil {
		t.Fatal(err)
	}
	if settings.SectionOrder[0] != "showcase" || settings.SectionOrder[2] != "catalog" {
		t.Fatal("legacy order was not preserved")
	}
}
