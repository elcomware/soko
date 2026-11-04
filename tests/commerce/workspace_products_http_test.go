package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDashboardProductsPaginationPreservesScopeAndSelectedRecord(t *testing.T) {
	owner := models.User{ID: "http-paged-owner", Email: "http-paged-owner@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&owner))
	for i, p := range []models.Product{{ID: "http-page-a", Title: "HTTP match a", Category: "One", VendorID: owner.ID}, {ID: "http-page-b", Title: "HTTP match b", Category: "Two", VendorID: owner.ID}, {ID: "http-page-hidden", Title: "HTTP match hidden", Category: "Secret", VendorID: "someone-else"}} {
		p.CreatedAt = time.Date(2026, 1, 1, 0, 0, i, 0, time.UTC)
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	tc := &frameworktesting.TestCase{}
	get := func(query string) map[string]any {
		response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard?section=products" + query)
		require.NoError(t, err)
		response.AssertOk()
		body, err := response.Json()
		require.NoError(t, err)
		return body["props"].(map[string]any)
	}
	props := get("&per_page=1&category=Two&selected_product=http-page-a")
	meta, ok := props["productPage"].(map[string]any)
	require.True(t, ok, "product section must expose authoritative pagination")
	require.EqualValues(t, 2, meta["total"])
	require.EqualValues(t, 1, meta["filtered_total"])
	require.Len(t, props["products"], 1)
	require.Equal(t, []any{"One", "Two"}, props["productCategories"])
	require.Equal(t, "http-page-a", props["selectedProduct"].(map[string]any)["id"])
	require.Equal(t, "available", props["selectedProductStatus"])
	require.True(t, props["recordCapabilities"].(map[string]any)["products"].(map[string]any)["http-page-a"].(map[string]any)["update"].(bool))
	for _, key := range []string{"orders", "offers", "rfqs", "users"} {
		require.Empty(t, props[key], key+" must not be loaded on products page")
	}
	props = get("&q=match%20a&selected_product=http-page-hidden")
	require.Equal(t, "unavailable", props["selectedProductStatus"])
	require.Nil(t, props["selectedProduct"])
	require.EqualValues(t, 1, props["productPage"].(map[string]any)["filtered_total"])
	_, leaked := props["recordCapabilities"].(map[string]any)["products"].(map[string]any)["http-page-hidden"]
	require.False(t, leaked)
	props = get("&q=missing&selected_product=http-page-a")
	require.Empty(t, props["products"])
	require.NotNil(t, props["selectedProduct"], "empty page must not hide authorized selected record")
	response, err := tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard?section=products&sort=unsafe")
	require.NoError(t, err)
	response.AssertStatus(400)
	// Removing permission must not leave an alternate paginated read path.
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: owner.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"}]`, Version: 1}))
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard?section=products")
	require.NoError(t, err)
	response.AssertForbidden()
}

func TestDashboardProductCountFailureDoesNotBecomeEmptySuccess(t *testing.T) {
	owner := models.User{ID: "http-paged-error", Email: "http-paged-error@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&owner))
	_, err := facades.Orm().Query().Exec("ALTER TABLE products RENAME TO http_unavailable_products")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE http_unavailable_products RENAME TO products")
		require.NoError(t, e)
	}()
	response, err := (&frameworktesting.TestCase{}).Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard?section=products")
	require.NoError(t, err)
	response.AssertStatus(500)
}
