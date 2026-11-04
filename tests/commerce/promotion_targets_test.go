package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPromotionTargetsScopeFacetsSelectionAndCoverage(t *testing.T) {
	actor := &models.User{ID: "pt-owner", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "promotions.manage", Scope: "own"}}}
	for _, p := range []models.Product{{ID: "pt-a", Title: "ÉCOLE", VendorID: actor.ID, Category: "A", SubCategory: "one"}, {ID: "pt-b", Title: "Beta", VendorID: actor.ID, Category: "A", SubCategory: "two"}, {ID: "pt-c", Title: "Gamma", VendorID: actor.ID, Category: "B", SubCategory: "one"}, {ID: "pt-hidden", Title: "Hidden", VendorID: "other", Category: "Secret"}} {
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	f := services.PromotionTargetFilters{VendorID: actor.ID, Query: "école", Page: 1, PerPage: 1, Scope: models.PromotionScope{All: true, ExcludeProductIDs: []string{"pt-b", "pt-hidden", "missing"}, ProductIDs: []string{"pt-c"}}}
	out, err := services.LoadPromotionTargets(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 3, out.ProductPage.Total)
	require.EqualValues(t, 1, out.ProductPage.FilteredTotal)
	require.Len(t, out.Products, 1)
	require.Equal(t, "pt-a", out.Products[0].ID)
	require.EqualValues(t, 2, out.Preview.Total)
	require.Len(t, out.Categories, 2)
	require.Len(t, out.Subcategories, 3)
	selected := map[string]string{}
	for _, p := range out.SelectedProducts {
		selected[p.ID] = p.Status
	}
	require.Equal(t, "available", selected["pt-c"])
	require.Equal(t, "unavailable", selected["pt-hidden"])
	require.Equal(t, "unavailable", selected["missing"])
	f.Query = ""
	f.Page = 999
	out, err = services.LoadPromotionTargets(actor, f)
	require.NoError(t, err)
	require.Equal(t, 3, out.ProductPage.Page)
	require.Equal(t, "pt-a", out.Products[0].ID)
	f.Scope = models.PromotionScope{Categories: []string{"A"}, ExcludeSubCategories: []models.CategorySelection{{Category: "A", SubCategory: "two"}}}
	out, err = services.LoadPromotionTargets(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, out.Preview.Total)
	f.VendorID = "other"
	out, err = services.LoadPromotionTargets(actor, f)
	require.Error(t, err)
	require.Nil(t, out)
	f.VendorID = ""
	_, err = services.LoadPromotionTargets(actor, f)
	require.Error(t, err)
}

func TestPromotionTargetsEmptySearchSelectionAndFailure(t *testing.T) {
	actor := &models.User{ID: "pt-edge", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "promotions.manage", Scope: "own"}}}
	for i, title := range []string{"ΟΣ", "ПРИВЕТ", "İSTANBUL", "100%_"} {
		p := models.Product{ID: fmt.Sprintf("pt-edge-%d", i), Title: title, VendorID: actor.ID, Category: " ", SubCategory: "orphan"}
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	for _, q := range []string{"ος", "привет", "i\u0307stanbul", "100%_"} {
		out, err := services.LoadPromotionTargets(actor, services.PromotionTargetFilters{VendorID: actor.ID, Query: q, Page: 1, PerPage: 1})
		require.NoError(t, err)
		require.EqualValues(t, 1, out.ProductPage.FilteredTotal)
		require.Zero(t, out.Preview.Total)
		require.Empty(t, out.Categories)
		require.Empty(t, out.Subcategories)
	}
	f := services.PromotionTargetFilters{VendorID: actor.ID, Query: "not present", Page: 999, PerPage: 1, Scope: models.PromotionScope{All: true, ProductIDs: []string{"pt-edge-0"}}}
	out, err := services.LoadPromotionTargets(actor, f)
	require.NoError(t, err)
	require.Empty(t, out.Products)
	require.EqualValues(t, 4, out.Preview.Total)
	require.Equal(t, "available", out.SelectedProducts[0].Status)
	f.Scope.ProductIDs = make([]string, 1001)
	require.Error(t, services.ValidatePromotionTargetFilters(f))
	f.Scope.ProductIDs = nil
	_, err = facades.Orm().Query().Exec("ALTER TABLE products RENAME TO pt_unavailable_products")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE pt_unavailable_products RENAME TO products")
		require.NoError(t, e)
	}()
	out, err = services.LoadPromotionTargets(actor, f)
	require.Error(t, err)
	require.Nil(t, out)
}
