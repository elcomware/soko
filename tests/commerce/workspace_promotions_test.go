package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWorkspacePromotionsScopedStatusesAndSelection(t *testing.T) {
	actor := &models.User{ID: "wp-owner", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "promotions.manage", Scope: "own"}}}
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "wp-product", VendorID: actor.ID, VendorName: "ÉCOLE Merchant", Title: "Only own target"}))
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "wp-other-product", VendorID: "foreign", Title: "Hidden target"}))
	for _, p := range []models.Promotion{{ID: "wp-a", Name: "Campaign active", VendorID: actor.ID, Enabled: true, StartsAt: now, EndsAt: now.Add(time.Hour)}, {ID: "wp-b", Name: "Campaign scheduled", VendorID: actor.ID, Enabled: true, StartsAt: now.Add(time.Minute), EndsAt: now.Add(time.Hour)}, {ID: "wp-c", Name: "Campaign ended", VendorID: actor.ID, Enabled: true, StartsAt: now.Add(-time.Hour), EndsAt: now}, {ID: "wp-d", Name: "Campaign disabled", VendorID: actor.ID, Enabled: false, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour)}, {ID: "wp-hidden", Name: "Hidden only", VendorID: "foreign", Enabled: true, StartsAt: now, EndsAt: now.Add(time.Hour)}} {
		p.Kind = "percentage"
		p.CreatedAt = now
		p.Scope = models.PromotionScope{All: true}
		p.PrepareSave()
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	f, err := services.ParseWorkspacePromotionFilters(map[string]string{"per_page": "1", "selected_campaign": "wp-a"})
	require.NoError(t, err)
	page, err := services.LoadWorkspacePromotions(actor, f, now)
	require.NoError(t, err)
	require.EqualValues(t, 4, page.CampaignPage.Total)
	require.Equal(t, "wp-d", page.Campaigns[0].ID)
	require.EqualValues(t, 1, page.CampaignCounts["active"])
	require.EqualValues(t, 1, page.CampaignCounts["scheduled"])
	require.EqualValues(t, 1, page.CampaignCounts["ended"])
	require.EqualValues(t, 1, page.CampaignCounts["disabled"])
	require.Equal(t, now.Add(time.Minute).Format(time.RFC3339Nano), page.NextStatusChange)
	require.True(t, page.SelectedCampaign.Scope.All)
	require.Len(t, page.Products, 1)
	f.Query = "école merchant"
	f.Status = "active"
	f.Page = 999
	f.SelectedCampaign = "wp-hidden"
	page, err = services.LoadWorkspacePromotions(actor, f, now)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.CampaignPage.FilteredTotal)
	require.Equal(t, 1, page.CampaignPage.Page)
	require.Nil(t, page.SelectedCampaign)
	require.Equal(t, "unavailable", page.SelectedCampaignStatus)
	f.Query = "Hidden only"
	f.Status = ""
	page, err = services.LoadWorkspacePromotions(actor, f, now)
	require.NoError(t, err)
	require.Zero(t, page.CampaignPage.FilteredTotal)
	f.Query = ""
	f.Status = "active"
	page, err = services.LoadWorkspacePromotions(actor, f, now.Add(time.Minute))
	require.NoError(t, err)
	require.EqualValues(t, 2, page.CampaignPage.FilteredTotal)
	f.Status = ""
	f.Sort = "oldest"
	f.Page = 1
	page, err = services.LoadWorkspacePromotions(actor, f, now)
	require.NoError(t, err)
	require.Equal(t, "wp-a", page.Campaigns[0].ID)
}

func TestWorkspacePromotionsUnicodeAndTargetFailure(t *testing.T) {
	actor := &models.User{ID: "wp-unicode", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "promotions.manage", Scope: "own"}}}
	now := time.Now().UTC()
	for i, name := range []string{"ΟΣ", "ПРИВЕТ", "İSTANBUL", "100%_"} {
		p := models.Promotion{ID: fmt.Sprintf("wp-u-%d", i), Name: name, VendorID: actor.ID, Kind: "fixed", Enabled: true, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour)}
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	for _, q := range []string{"ος", "привет", "i\u0307stanbul", "100%_"} {
		f, err := services.ParseWorkspacePromotionFilters(map[string]string{"q": q})
		require.NoError(t, err)
		page, err := services.LoadWorkspacePromotions(actor, f, now)
		require.NoError(t, err)
		require.EqualValues(t, 1, page.CampaignPage.FilteredTotal, q)
	}
	f, err := services.ParseWorkspacePromotionFilters(map[string]string{"selected_campaign": "wp-u-0"})
	require.NoError(t, err)
	_, err = facades.Orm().Query().Exec("ALTER TABLE products RENAME TO wp_unavailable_products")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE wp_unavailable_products RENAME TO products")
		require.NoError(t, e)
	}()
	page, err := services.LoadWorkspacePromotions(actor, f, now)
	require.Error(t, err)
	require.Nil(t, page)
}
func TestWorkspacePromotionFiltersRejectInvalid(t *testing.T) {
	for _, input := range []map[string]string{{"status": "pending"}, {"sort": "price"}, {"page": "0"}, {"per_page": "101"}, {"q": "bad\nquery"}} {
		_, err := services.ParseWorkspacePromotionFilters(input)
		require.Error(t, err)
	}
	f, err := services.ParseWorkspacePromotionFilters(nil)
	require.NoError(t, err)
	_, err = services.LoadWorkspacePromotions(nil, f, time.Now())
	require.Error(t, err)
}
