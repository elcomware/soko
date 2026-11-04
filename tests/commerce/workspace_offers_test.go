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

func TestWorkspaceOffersScopeSelectionAndOrdering(t *testing.T) {
	actor := &models.User{ID: "wof-owner", Role: models.RoleMerchant}
	stamp := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, p := range []models.Product{{ID: "wof-product", VendorID: actor.ID}, {ID: "wof-other", VendorID: "foreign"}} {
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	for _, o := range []models.Offer{{ID: "wof-a", ProductID: "wof-product", ProductTitle: "ÉCOLE", Status: models.OfferStatusPending, CreatedAt: stamp}, {ID: "wof-b", ProductID: "wof-product", Status: models.OfferStatusAccepted, CreatedAt: stamp}, {ID: "wof-hidden", ProductID: "wof-other", ProductTitle: "Hidden only", Status: models.OfferStatusPending, CreatedAt: stamp}, {ID: "wof-orphan", ProductID: "missing", Status: models.OfferStatusPending}} {
		require.NoError(t, facades.Orm().Query().Create(&o))
	}
	f, err := services.ParseWorkspaceOfferFilters(map[string]string{"per_page": "1", "selected_offer": "wof-a"})
	require.NoError(t, err)
	page, err := services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 2, page.OfferPage.Total)
	require.Equal(t, "wof-b", page.Offers[0].ID)
	require.Equal(t, "supplier", page.OfferVisibility)
	require.True(t, page.RecordCapabilities["offers"]["wof-a"]["decide"])
	f.Query = "école"
	f.Page = 999
	f.SelectedOffer = "wof-hidden"
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.OfferPage.FilteredTotal)
	require.Equal(t, 1, page.OfferPage.Page)
	require.Nil(t, page.SelectedOffer)
	require.Equal(t, "unavailable", page.SelectedOfferStatus)
	f.Query = "Hidden only"
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.Zero(t, page.OfferPage.FilteredTotal)
	f.Query = ""
	f.Status = "accepted"
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.OfferPage.FilteredTotal)
	require.False(t, page.RecordCapabilities["offers"]["wof-b"]["decide"])
	f.Status = ""
	f.Sort = "oldest"
	f.Page = 1
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.Equal(t, "wof-a", page.Offers[0].ID)
	f.SelectedOffer = "wof-orphan"
	page, err = services.LoadWorkspaceOffers(&models.User{ID: "wof-admin", Role: models.RoleAdmin}, f)
	require.NoError(t, err)
	require.Equal(t, "full", page.OfferVisibility)
	require.NotNil(t, page.SelectedOffer)
	require.False(t, page.RecordCapabilities["offers"]["wof-orphan"]["decide"])
	// A product transfer removes supplier read and decision access on refresh.
	_, err = facades.Orm().Query().Model(&models.Product{}).Where("id = ?", "wof-product").Update("vendor_id", "new-vendor")
	require.NoError(t, err)
	f.SelectedOffer = "wof-a"
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.Zero(t, page.OfferPage.Total)
	require.Nil(t, page.SelectedOffer)
	require.Empty(t, page.RecordCapabilities["offers"])
}

func TestWorkspaceOffersUnicodeAndReadAllWriteOwn(t *testing.T) {
	actor := &models.User{ID: "wof-mixed", Role: models.RoleMerchant, PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "workspace.access", Scope: "own"}, {Permission: "offers.read", Scope: "all"}, {Permission: "offers.decide", Scope: "own"}}}
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"},{"permission":"offers.read","scope":"all"},{"permission":"offers.decide","scope":"own"}]`, Version: 1}))
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "wof-mixed-product", VendorID: actor.ID}))
	for i, title := range []string{"wof-u ΟΣ", "wof-u ПРИВЕТ", "wof-u İSTANBUL", "wof-u 100%_"} {
		id := fmt.Sprintf("wof-u-%d", i)
		require.NoError(t, facades.Orm().Query().Create(&models.Offer{ID: id, ProductID: "wof-mixed-product", ProductTitle: title, Status: models.OfferStatusPending}))
	}
	for _, q := range []string{"wof-u ος", "wof-u привет", "wof-u i\u0307stanbul", "wof-u 100%_"} {
		f, err := services.ParseWorkspaceOfferFilters(map[string]string{"q": q})
		require.NoError(t, err)
		page, err := services.LoadWorkspaceOffers(actor, f)
		require.NoError(t, err)
		require.EqualValues(t, 1, page.OfferPage.FilteredTotal, q)
		require.True(t, page.RecordCapabilities["offers"][page.Offers[0].ID]["decide"])
	}
	f, err := services.ParseWorkspaceOfferFilters(map[string]string{"selected_offer": "wof-u-0"})
	require.NoError(t, err)
	require.NoError(t, services.DecideOffer("wof-u-0", "accepted", actor))
	require.Error(t, services.DecideOffer("wof-u-0", "declined", actor))
	page, err := services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.False(t, page.RecordCapabilities["offers"]["wof-u-0"]["decide"])
	_, err = facades.Orm().Query().Model(&models.Product{}).Where("id = ?", "wof-mixed-product").Update("vendor_id", "someone-else")
	require.NoError(t, err)
	f.SelectedOffer = "wof-u-1"
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.NoError(t, err)
	require.NotNil(t, page.SelectedOffer)
	require.False(t, page.RecordCapabilities["offers"]["wof-u-1"]["decide"])
	// Owner query errors must not become permissive capabilities or empty success.
	_, err = facades.Orm().Query().Exec("ALTER TABLE products RENAME TO wof_unavailable_products")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE wof_unavailable_products RENAME TO products")
		require.NoError(t, e)
	}()
	page, err = services.LoadWorkspaceOffers(actor, f)
	require.Error(t, err)
	require.Nil(t, page)
}
func TestWorkspaceOfferFilterValidation(t *testing.T) {
	for _, input := range []map[string]string{{"status": "open"}, {"sort": "price"}, {"page": "0"}, {"per_page": "101"}, {"q": "bad\nquery"}} {
		_, err := services.ParseWorkspaceOfferFilters(input)
		require.Error(t, err)
	}
	f, err := services.ParseWorkspaceOfferFilters(nil)
	require.NoError(t, err)
	_, err = services.LoadWorkspaceOffers(nil, f)
	require.Error(t, err)
}
