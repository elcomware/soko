package models

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCampaignScopeWindowAndExclusions(t *testing.T) {
	now := time.Date(2026, 12, 20, 12, 0, 0, 0, time.UTC)
	p := Product{ID: "coat", Category: "Clothing", SubCategory: "Coats", VendorID: "seller", Price: 500, MOQ: 1, Stock: 10}
	c := Promotion{ID: "xmas", Name: "Christmas", Enabled: true, StartsAt: now, EndsAt: now.Add(time.Hour), Kind: "percentage", Amount: 20, Scope: PromotionScope{Categories: []string{"Clothing"}, ExcludeProductIDs: []string{"hat"}}}
	require.NoError(t, c.Validate())
	require.True(t, c.Applies(p, now))
	require.False(t, c.Applies(p, now.Add(-time.Nanosecond)))
	require.False(t, c.Applies(p, c.EndsAt))
	c.Scope.ExcludeSubCategories = []CategorySelection{{Category: "Clothing", SubCategory: "Coats"}}
	require.False(t, c.Applies(p, now))
	c.Scope = PromotionScope{All: true}
	c.VendorID = "other"
	require.False(t, c.Applies(p, now))
}
func TestPriceCampaignBestDiscountAndTierSnapshot(t *testing.T) {
	now := time.Now()
	p := Product{ID: "a", Category: "Tools", Price: 500, MOQ: 1, Stock: 20, PriceTiers: []PriceTier{{MinQty: 1, MaxQty: 2, Price: 500}, {MinQty: 3, Price: 400}}}
	base := Promotion{ID: "percent", Name: "Sale", Enabled: true, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour), Kind: "percentage", Amount: 20, Scope: PromotionScope{All: true}}
	fixed := base
	fixed.ID = "fixed"
	fixed.Kind = "fixed"
	fixed.Amount = 100
	price, err := PriceProduct(&p, 3, []Promotion{base, fixed}, now)
	require.NoError(t, err)
	require.Equal(t, 400.0, price.OriginalUnitPrice)
	require.Equal(t, 300.0, price.UnitPrice)
	require.Equal(t, "fixed", price.PromotionID)
	expired, err := PriceProduct(&p, 3, []Promotion{base}, base.EndsAt)
	require.NoError(t, err)
	require.Equal(t, 400.0, expired.UnitPrice)
	require.Empty(t, expired.PromotionID)
	fixed.Amount = 10000
	price, err = PriceProduct(&p, 3, []Promotion{fixed}, now)
	require.NoError(t, err)
	require.Equal(t, 0.0, price.UnitPrice)
}
func TestCampaignValidation(t *testing.T) {
	now := time.Now()
	c := Promotion{Name: "Sale", Kind: "percentage", Amount: 20, StartsAt: now, EndsAt: now.Add(time.Hour), Scope: PromotionScope{All: true}}
	require.NoError(t, c.Validate())
	c.Amount = 101
	require.Error(t, c.Validate())
	c.Amount = 10
	c.EndsAt = now
	require.Error(t, c.Validate())
	c.EndsAt = now.Add(time.Hour)
	c.Scope = PromotionScope{}
	require.Error(t, c.Validate())
}

func TestCampaignRoundingTieBreakAndNoFabricatedDiscount(t *testing.T) {
	now := time.Now()
	p := Product{ID: "round", Price: 19.99, MOQ: 1, Stock: 10}
	c := Promotion{ID: "b", Name: "Ten percent", Enabled: true, Kind: "percentage", Amount: 10, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute), Scope: PromotionScope{ProductIDs: []string{p.ID}}}
	same := c
	same.ID = "a"
	price, err := PriceProduct(&p, 1, []Promotion{c, same}, now)
	require.NoError(t, err)
	require.Equal(t, 17.99, price.UnitPrice)
	require.Equal(t, 2.0, price.DiscountPerUnit)
	require.Equal(t, "a", price.PromotionID)
	same.Enabled = false
	c.Scope.ExcludeCategories = []string{p.Category}
	price, err = PriceProduct(&p, 1, []Promotion{c, same}, now)
	require.NoError(t, err)
	require.Equal(t, 19.99, price.UnitPrice)
	require.Empty(t, price.PromotionID)
	require.Zero(t, price.DiscountPerUnit)
}
func TestCampaignSubcategoryRequiresCategoryMatch(t *testing.T) {
	now := time.Now()
	c := Promotion{Enabled: true, StartsAt: now, EndsAt: now.Add(time.Hour), Scope: PromotionScope{SubCategories: []CategorySelection{{Category: "Clothing", SubCategory: "Accessories"}}}}
	require.True(t, c.Applies(Product{Category: "Clothing", SubCategory: "Accessories"}, now))
	require.False(t, c.Applies(Product{Category: "Electronics", SubCategory: "Accessories"}, now))
}
