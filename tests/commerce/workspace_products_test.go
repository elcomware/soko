package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWorkspaceProductPageScopesCountsAndStableOrdering(t *testing.T) {
	actor := &models.User{ID: "paged-products-owner", Role: models.RoleMerchant}
	stamp := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, p := range []models.Product{
		{ID: "paged-a", VendorID: actor.ID, Title: "Paged needle", Category: "Paged", CreatedAt: stamp, VendorRegion: "East Africa", GalleryUrlsJSON: `["https://example.com/a.png"]`},
		{ID: "paged-b", VendorID: actor.ID, Title: "Paged 100%_ literal", Category: "Paged", CreatedAt: stamp},
		{ID: "paged-c", VendorID: actor.ID, Title: "Paged newer", Category: "Other", CreatedAt: stamp.Add(time.Hour)},
		{ID: "paged-foreign", VendorID: "another-owner", Title: "Paged needle", Category: "Paged", CreatedAt: stamp.Add(2 * time.Hour)},
	} {
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	require.NoError(t, facades.Orm().Query().Create(&models.Offer{ID: "paged-history", ProductID: "paged-b"}))
	filters, err := services.ParseWorkspaceProductFilters(map[string]string{"per_page": "1"})
	require.NoError(t, err)
	first, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.EqualValues(t, 3, first.ProductPage.Total)
	require.EqualValues(t, 3, first.ProductPage.FilteredTotal)
	require.Equal(t, 3, first.ProductPage.LastPage)
	require.Len(t, first.Products, 1)
	require.Equal(t, "paged-c", first.Products[0].ID)
	require.Equal(t, 1, first.ProductPage.From)
	require.Equal(t, 1, first.ProductPage.To)
	filters.Page = 2
	second, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.Equal(t, "paged-b", second.Products[0].ID)
	require.False(t, second.RecordCapabilities["products"]["paged-b"]["delete"])
	require.True(t, second.RecordCapabilities["products"]["paged-b"]["update"])
	filters.Page = 3
	third, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.Equal(t, "paged-a", third.Products[0].ID)
	require.Equal(t, "east-africa", third.Products[0].VendorRegion)
	require.Len(t, third.Products[0].GalleryUrls, 1)
	require.True(t, third.RecordCapabilities["products"]["paged-a"]["delete"])
	filters, err = services.ParseWorkspaceProductFilters(map[string]string{"q": "needle", "category": "Paged", "page": "999"})
	require.NoError(t, err)
	found, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.EqualValues(t, 3, found.ProductPage.Total)
	require.EqualValues(t, 1, found.ProductPage.FilteredTotal)
	require.Equal(t, 1, found.ProductPage.Page)
	require.Equal(t, "paged-a", found.Products[0].ID)
	filters.Query = "100%_"
	literal, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.Len(t, literal.Products, 1)
	require.Equal(t, "paged-b", literal.Products[0].ID)
	filters.Query = "absent"
	empty, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.NotNil(t, empty.Products)
	require.Empty(t, empty.Products)
	require.EqualValues(t, 0, empty.ProductPage.FilteredTotal)
	require.Zero(t, empty.ProductPage.From)
	require.Zero(t, empty.ProductPage.To)
	all := &models.User{ID: "paged-reader", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "products.read", Scope: "all"}}}
	filters.Query = "Paged"
	filters.Category = ""
	filters.Sort = "oldest"
	filters.PerPage = 25
	broad, err := services.LoadWorkspaceProducts(all, filters)
	require.NoError(t, err)
	require.Len(t, broad.Products, 4)
	require.Equal(t, "paged-a", broad.Products[0].ID)
	require.False(t, broad.RecordCapabilities["products"]["paged-a"]["update"])
	filters.Sort = "title_asc"
	titled, err := services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.Equal(t, []string{"paged-b", "paged-a", "paged-c"}, []string{titled.Products[0].ID, titled.Products[1].ID, titled.Products[2].ID})
	filters.Sort = "title_desc"
	titled, err = services.LoadWorkspaceProducts(actor, filters)
	require.NoError(t, err)
	require.Equal(t, []string{"paged-c", "paged-a", "paged-b"}, []string{titled.Products[0].ID, titled.Products[1].ID, titled.Products[2].ID})

	revoked := &models.User{ID: actor.ID, Role: models.RoleAdmin, PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{}}
	_, err = services.LoadWorkspaceProducts(revoked, filters)
	require.Error(t, err)
}

func TestWorkspaceProductFiltersRejectUnsafeAndUnboundedInputs(t *testing.T) {
	defaults, err := services.ParseWorkspaceProductFilters(nil)
	require.NoError(t, err)
	require.Equal(t, 1, defaults.Page)
	require.Equal(t, 25, defaults.PerPage)
	require.Equal(t, "newest", defaults.Sort)
	for _, input := range []map[string]string{{"page": "0"}, {"page": "-1"}, {"page": "1.5"}, {"page": "1000001"}, {"per_page": "101"}, {"per_page": "0"}, {"sort": "title; DROP TABLE products"}, {"q": string(make([]byte, 121))}, {"category": string(make([]byte, 129))}} {
		_, err := services.ParseWorkspaceProductFilters(input)
		require.Error(t, err, input)
	}
	_, err = services.LoadWorkspaceProducts(&models.User{ID: "unsafe-page", Role: models.RoleAdmin}, services.WorkspaceProductFilters{Page: -1, PerPage: 25, Sort: "newest"})
	require.Error(t, err)
}

func TestWorkspaceProductsDoNotLoadOtherCollectionsAndPropagateStorageErrors(t *testing.T) {
	actor := &models.User{ID: "isolated-products", Role: models.RoleMerchant}
	filters, err := services.ParseWorkspaceProductFilters(nil)
	require.NoError(t, err)
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "isolated-product", VendorID: actor.ID}))
	// The page must not depend on users/RFQs merely because dashboard normally loads them.
	_, err = facades.Orm().Query().Exec("ALTER TABLE rfqs RENAME TO unavailable_paged_rfqs")
	require.NoError(t, err)
	func() {
		defer func() {
			_, e := facades.Orm().Query().Exec("ALTER TABLE unavailable_paged_rfqs RENAME TO rfqs")
			require.NoError(t, e)
		}()
		page, e := services.LoadWorkspaceProducts(actor, filters)
		require.NoError(t, e)
		require.Len(t, page.Products, 1)
	}()
	_, err = facades.Orm().Query().Exec("ALTER TABLE products RENAME TO unavailable_paged_products")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE unavailable_paged_products RENAME TO products")
		require.NoError(t, e)
	}()
	page, err := services.LoadWorkspaceProducts(actor, filters)
	require.Error(t, err)
	require.Nil(t, page)
}

func TestWorkspaceProductUnicodeSearchKeepsCountsScopeAndPaging(t *testing.T) {
	actor := &models.User{ID: "unicode-paged-owner", Role: models.RoleMerchant}
	stamp := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	for i, p := range []models.Product{
		{ID: "unicode-a", Title: "ÉCOLE Kit", VendorID: actor.ID, Category: "Learning"},
		{ID: "unicode-b", Title: "école Bag", VendorID: actor.ID, Category: "Learning"},
		{ID: "unicode-c", Title: "ÉCOLE Tool", VendorID: actor.ID, Category: "Tools"},
		{ID: "unicode-other", Title: "ÉCOLE Private", VendorID: "foreign", Category: "Learning"},
		{ID: "unicode-greek", Title: "ΟΣ", VendorID: actor.ID, Category: "Languages"},
		{ID: "unicode-turkish", VendorName: "İSTANBUL", VendorID: actor.ID, Category: "Languages"},
		{ID: "unicode-cyrillic", Title: "МОСКВА", VendorID: actor.ID, Category: "Languages"},
	} {
		p.CreatedAt = stamp.Add(time.Duration(i) * time.Second)
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	f, err := services.ParseWorkspaceProductFilters(map[string]string{"q": "école", "category": "Learning", "per_page": "1", "page": "99", "sort": "oldest"})
	require.NoError(t, err)
	result, err := services.LoadWorkspaceProducts(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 6, result.ProductPage.Total)
	require.EqualValues(t, 2, result.ProductPage.FilteredTotal)
	require.Equal(t, 2, result.ProductPage.Page)
	require.Equal(t, "unicode-b", result.Products[0].ID)
	for _, fixture := range []struct{ query, id string }{{"ΟΣ", "unicode-greek"}, {"i\u0307stanbul", "unicode-turkish"}, {"москва", "unicode-cyrillic"}} {
		f.Query = fixture.query
		f.Category = ""
		f.Page = 1
		result, err = services.LoadWorkspaceProducts(actor, f)
		require.NoError(t, err)
		require.Len(t, result.Products, 1, fixture.query)
		require.Equal(t, fixture.id, result.Products[0].ID)
	}
}

func TestWorkspaceProductSearchCursorErrorsPropagate(t *testing.T) {
	actor := &models.User{ID: "cursor-error-owner", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "cursor-error-product", VendorID: actor.ID, Title: "Search me"}))
	f, err := services.ParseWorkspaceProductFilters(map[string]string{"q": "Search"})
	require.NoError(t, err)
	_, err = facades.Orm().Query().Exec("ALTER TABLE products RENAME COLUMN title TO unavailable_search_title")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE products RENAME COLUMN unavailable_search_title TO title")
		require.NoError(t, e)
	}()
	// Counts still succeed, but a failed search cursor must not claim zero matches.
	result, err := services.LoadWorkspaceProducts(actor, f)
	require.Error(t, err)
	require.Nil(t, result)
}
