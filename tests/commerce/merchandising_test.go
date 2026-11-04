package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
	"elcomware/app/services"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMerchandisingPersistsWithScopeAndConcurrency(t *testing.T) {
	c := queries.DefaultMerchandising()
	admin := &models.User{ID: "merch-admin", Role: models.RoleAdmin}
	merchant := &models.User{ID: "merch-supplier", Role: models.RoleMerchant}
	require.Error(t, services.SaveMerchandising(merchant, c))
	c.Hero.Title = "Real selected products"
	c.SectionOrder = []string{"catalog", "showcase", "discovery", "category_navigation", "sourcing"}
	require.NoError(t, services.SaveMerchandising(admin, c))
	stored, err := queries.LoadMerchandising()
	require.NoError(t, err)
	require.Equal(t, 1, stored.Version)
	require.Equal(t, c.Hero.Title, stored.Hero.Title)
	require.Equal(t, c.SectionOrder, stored.SectionOrder)
	c.Hero.Title = "Stale overwrite"
	require.Error(t, services.SaveMerchandising(admin, c))
	stored.Hero.Strategy = "manual"
	stored.Hero.ProductIDs = []string{"nonexistent-product"}
	require.Error(t, services.SaveMerchandising(admin, stored))
	saved, err := queries.LoadMerchandising()
	require.NoError(t, err)
	require.Equal(t, "Real selected products", saved.Hero.Title)
	saved.Hero.ProductIDs = []string{"", "deleted-pin"}
	saved.Hero.Strategy = "newest"
	require.NoError(t, services.SaveMerchandising(admin, saved))
	cleared, err := queries.LoadMerchandising()
	require.NoError(t, err)
	require.Empty(t, cleared.Hero.ProductIDs)
}
func TestMerchandisingHTTPRequiresCapability(t *testing.T) {
	u := models.User{ID: "merch-no-access", Email: "merch-no-access@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&u))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithSession(map[string]any{"user_id": u.ID}).Get("/dashboard/merchandising")
	require.NoError(t, err)
	response.AssertForbidden()
	response, err = tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"}).WithSession(map[string]any{"user_id": u.ID}).Post("/dashboard/merchandising", strings.NewReader(`{}`))
	require.NoError(t, err)
	response.AssertForbidden()
}

func TestStorefrontCategoryChoicesIncludeNewCategoriesOutsideCurrentFilter(t *testing.T) {
	for _, p := range []models.Product{{ID: "merch-category-a", Title: "UniqueScopeItem", Category: "New Category One", Type: models.TypePhysical, Price: 10, Stock: 5, MOQ: 1}, {ID: "merch-category-b", Title: "Other item", Category: "New Category Two", Type: models.TypePhysical, Price: 12, Stock: 5, MOQ: 1}} {
		product := p
		require.NoError(t, facades.Orm().Query().Create(&product))
	}
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeader("X-Inertia", "true").Get("/?q=UniqueScopeItem")
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	categories := props["categories"].([]any)
	found := false
	for _, v := range categories {
		if v.(map[string]any)["id"] == "New Category Two" {
			found = true
		}
	}
	require.True(t, found, "new categories should remain navigable outside the active filter")
	products := props["products"].([]any)
	require.Len(t, products, 1)
	require.Equal(t, "merch-category-a", products[0].(map[string]any)["id"])
	merch := props["merchandising"].(map[string]any)
	for _, key := range []string{"hero", "hero_side", "discovery"} {
		for _, p := range merch[key].(map[string]any)["products"].([]any) {
			require.Equal(t, "merch-category-a", p.(map[string]any)["id"], "placement escaped filter")
		}
	}
}
