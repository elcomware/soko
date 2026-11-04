package commerce_test

import (
	"bytes"
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/services"
	"encoding/base64"
	"encoding/json"
	"fmt"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"os"
	"strings"
	"testing"
)

func TestMediaHTTPAuthorizationValidationAndPersistence(t *testing.T) {
	buyer := models.User{ID: "media-buyer", Email: "media-buyer@example.test", Role: models.RoleBuyer}
	owner := models.User{ID: "media-owner", Email: "media-owner@example.test", Role: models.RoleMerchant}
	other := models.User{ID: "media-other", Email: "media-other@example.test", Role: models.RoleMerchant}
	for _, u := range []*models.User{&buyer, &owner, &other} {
		require.NoError(t, facades.Orm().Query().Create(u))
	}
	tc := &frameworktesting.TestCase{}
	resp, err := tc.Http(t).WithSession(map[string]any{"user_id": buyer.ID}).Get("/dashboard/media")
	require.NoError(t, err)
	resp.AssertForbidden()
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard/media")
	require.NoError(t, err)
	resp.AssertOk()
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 800, 600))))
	payload, _ := json.Marshal(map[string]string{"data": "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())})
	resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": owner.ID}).Post("/dashboard/media/upload", strings.NewReader(string(payload)))
	require.NoError(t, err)
	resp.AssertCreated()
	body, err := resp.Json()
	require.NoError(t, err)
	asset := body["asset"].(map[string]any)
	path := "storage/app/product-media/" + asset["id"].(string) + ".png"
	t.Cleanup(func() { _ = os.Remove(path) })
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, "png", format)
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": other.ID}).Get("/dashboard/media")
	require.NoError(t, err)
	body, err = resp.Json()
	require.NoError(t, err)
	require.Empty(t, body["assets"])
	policy, _ := json.Marshal(services.DefaultMediaPolicy())
	resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": owner.ID}).Post("/dashboard/media/policy", strings.NewReader(string(policy)))
	require.NoError(t, err)
	resp.AssertForbidden()
	bad := `{"template":"software","title":"Tool","brand":"Acme","primary":"url(x)","accent":"#123456","background":"#ffffff"}`
	resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": owner.ID}).Post("/dashboard/media/mockup", strings.NewReader(bad))
	require.NoError(t, err)
	resp.AssertUnprocessableEntity()
}
func TestMediaMockupAndAdminPolicy(t *testing.T) {
	admin := models.User{ID: "media-admin", Email: "media-admin@example.test", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&admin))
	tc := &frameworktesting.TestCase{}
	opts := services.MockupOptions{Template: "box", Title: "A useful product", Brand: "My brand", Primary: "#14283d", Accent: "#d5a64d", Background: "#f4f5f7"}
	b, _ := json.Marshal(opts)
	resp, err := tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/media/mockup", strings.NewReader(string(b)))
	require.NoError(t, err)
	resp.AssertCreated()
	body, err := resp.Json()
	require.NoError(t, err)
	asset := body["asset"].(map[string]any)
	file := "storage/app/product-media/" + asset["id"].(string) + ".svg"
	t.Cleanup(func() { _ = os.Remove(file) })
	svg, err := os.ReadFile(file)
	require.NoError(t, err)
	require.Contains(t, string(svg), "Illustrative mockup")
	require.Contains(t, string(svg), "My brand")
	policy := services.DefaultMediaPolicy()
	policy.AllowMockups = false
	encoded, _ := json.Marshal(policy)
	resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/media/policy", strings.NewReader(string(encoded)))
	require.NoError(t, err)
	resp.AssertOk()
	require.False(t, services.LoadMediaPolicy().AllowMockups)
	t.Cleanup(func() {
		_, _ = facades.Orm().Query().Where("id = ?", "media-policy").Delete(&models.StorefrontContent{})
	})
	resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/media/mockup", strings.NewReader(string(b)))
	require.NoError(t, err)
	resp.AssertUnprocessableEntity()
}
func TestMediaDraftIsPrivateUntilOwnedListingPublication(t *testing.T) {
	owner := models.User{ID: "media-private-owner", Email: "media-private-owner@example.test", Role: models.RoleMerchant}
	other := models.User{ID: "media-private-other", Email: "media-private-other@example.test", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&owner))
	require.NoError(t, facades.Orm().Query().Create(&other))
	tc := &frameworktesting.TestCase{}
	b, _ := json.Marshal(services.MockupOptions{Template: "book", Title: "Private draft", Brand: "Brand", Primary: "#14283d", Accent: "#d5a64d", Background: "#f4f5f7"})
	resp, err := tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": owner.ID}).Post("/dashboard/media/mockup", strings.NewReader(string(b)))
	require.NoError(t, err)
	resp.AssertCreated()
	body, err := resp.Json()
	require.NoError(t, err)
	asset := body["asset"].(map[string]any)
	assetURL := asset["url"].(string)
	resp, err = tc.Http(t).Get(assetURL)
	require.NoError(t, err)
	resp.AssertNotFound()
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": other.ID}).Get(assetURL)
	require.NoError(t, err)
	resp.AssertNotFound()
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get(assetURL)
	require.NoError(t, err)
	resp.AssertOk()
	input := &requests.CreateProductRequest{Title: "Private picture product", Category: "Books", Type: "physical", Format: "retail", Condition: "new", Price: 10, MOQ: 1, Stock: 1, ImageUrl: assetURL}
	_, err = services.NewCatalogService().CreateProduct(input, &other)
	require.Error(t, err)
	_, err = services.NewCatalogService().CreateProduct(input, &owner)
	require.NoError(t, err)
	resp, err = tc.Http(t).Get(assetURL)
	require.NoError(t, err)
	resp.AssertOk()
	t.Cleanup(func() { _ = os.Remove("storage/app/product-media/" + asset["id"].(string) + ".svg") })
}
func TestDisabledRemoteImagePolicyRejectsMixedCaseScheme(t *testing.T) {
	p := services.DefaultMediaPolicy()
	p.AllowRemote = false
	b, _ := json.Marshal(p)
	var row models.StorefrontContent
	require.NoError(t, facades.Orm().Query().UpdateOrCreate(&row, models.StorefrontContent{ID: "media-policy"}, map[string]any{"content_json": string(b)}))
	t.Cleanup(func() {
		_, _ = facades.Orm().Query().Where("id = ?", "media-policy").Delete(&models.StorefrontContent{})
	})
	input := &requests.CreateProductRequest{Title: "Image product", Category: "Books", Type: "physical", Format: "retail", Condition: "new", Price: 10, MOQ: 1, Stock: 1, ImageUrl: "HTTPS://example.com/a.png"}
	_, err := services.NewCatalogService().CreateProduct(input, &models.User{ID: "media-policy-merchant", Role: models.RoleMerchant})
	require.ErrorContains(t, err, "remote image")
}
func TestExternalMediaPathDoesNotRequireLocalAssetOwnership(t *testing.T) {
	input := &requests.CreateProductRequest{Title: "External image product", Category: "Books", Type: "physical", Format: "retail", Condition: "new", Price: 10, MOQ: 1, Stock: 1, ImageUrl: "https://cdn.example.com/media/photo.png"}
	_, err := services.NewCatalogService().CreateProduct(input, &models.User{ID: "media-external-merchant", Role: models.RoleMerchant})
	require.NoError(t, err)
}

func TestMediaLibraryPaginationSearchAndIsolation(t *testing.T) {
	owner := models.User{ID: "media-pages-owner", Email: "media-pages@example.test", Role: models.RoleMerchant}
	other := models.User{ID: "media-pages-other", Email: "media-pages-other@example.test", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&owner))
	require.NoError(t, facades.Orm().Query().Create(&other))
	for i := 0; i < 27; i++ {
		asset := models.ProductMedia{ID: fmt.Sprintf("media-page-%02d", i), OwnerID: owner.ID, URL: fmt.Sprintf("/media/page-%d", i), Kind: "mockup", OptionsJSON: `{"title":"Coffee packaging"}`}
		require.NoError(t, facades.Orm().Query().Create(&asset))
	}
	asset := models.ProductMedia{ID: "media-hidden-page", OwnerID: other.ID, URL: "/media/hidden", Kind: "mockup", OptionsJSON: `{"title":"Coffee packaging"}`}
	require.NoError(t, facades.Orm().Query().Create(&asset))
	tc := &frameworktesting.TestCase{}
	resp, err := tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard/media?page=2&q=Coffee&kind=mockup")
	require.NoError(t, err)
	resp.AssertOk()
	body, err := resp.Json()
	require.NoError(t, err)
	require.Len(t, body["assets"], 3)
	require.EqualValues(t, 27, body["total"])
	require.EqualValues(t, 2, body["page"])
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard/media?q=unmatched")
	require.NoError(t, err)
	body, err = resp.Json()
	require.NoError(t, err)
	require.Empty(t, body["assets"])
}

func TestMediaWorkspaceSupportsScopedDelegates(t *testing.T) {
	tc := &frameworktesting.TestCase{}
	for _, scenario := range []struct {
		id, permission, scope string
		create, policy        bool
	}{
		{"media-only-delegate", "media.create", "own", true, false},
		{"media-policy-delegate", "media.settings.manage", "all", false, true},
		{"media-denied-delegate", "products.view", "own", false, false},
	} {
		t.Run(scenario.id, func(t *testing.T) {
			u := models.User{ID: scenario.id, Email: scenario.id + "@example.test", Role: models.RoleBuyer}
			require.NoError(t, facades.Orm().Query().Create(&u))
			grants, _ := json.Marshal([]models.PermissionGrant{{Permission: scenario.permission, Scope: scenario.scope}})
			require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: u.ID, GrantsJSON: string(grants), Version: 1}))
			resp, err := tc.Http(t).WithSession(map[string]any{"user_id": u.ID}).Get("/dashboard/media-library")
			require.NoError(t, err)
			if !scenario.create && !scenario.policy {
				resp.AssertForbidden()
				return
			}
			resp.AssertOk()
			resp, err = tc.Http(t).WithSession(map[string]any{"user_id": u.ID}).Get("/dashboard/media")
			require.NoError(t, err)
			resp.AssertOk()
			body, err := resp.Json()
			require.NoError(t, err)
			require.Equal(t, scenario.create, body["can_create"])
			require.Equal(t, scenario.policy, body["can_manage_policy"])
			if !scenario.create {
				resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": u.ID}).Post("/dashboard/media/upload", strings.NewReader(`{"data":"bad"}`))
				require.NoError(t, err)
				resp.AssertForbidden()
			}
		})
	}
}

func TestMediaMockupRecipeIsOwnerScoped(t *testing.T) {
	owner := models.User{ID: "recipe-owner", Email: "recipe-owner@example.test", Role: models.RoleMerchant}
	other := models.User{ID: "recipe-other", Email: "recipe-other@example.test", Role: models.RoleMerchant}
	buyer := models.User{ID: "recipe-buyer", Email: "recipe-buyer@example.test", Role: models.RoleBuyer}
	for _, u := range []*models.User{&owner, &other, &buyer} {
		require.NoError(t, facades.Orm().Query().Create(u))
	}
	asset := models.ProductMedia{ID: "recipe-mockup", OwnerID: owner.ID, URL: "/media/recipe-mockup", Kind: "mockup", OptionsJSON: `{"template":"box","title":"Coffee","brand":"My brand","primary":"#123456","accent":"#654321","background":"#ffffff","logo_data":"must-not-leak"}`}
	photo := models.ProductMedia{ID: "recipe-photo", OwnerID: owner.ID, URL: "/media/recipe-photo", Kind: "upload"}
	require.NoError(t, facades.Orm().Query().Create(&asset))
	require.NoError(t, facades.Orm().Query().Create(&photo))
	tc := &frameworktesting.TestCase{}
	resp, err := tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard/media/recipe-mockup/recipe")
	require.NoError(t, err)
	resp.AssertOk()
	body, err := resp.Json()
	require.NoError(t, err)
	recipe := body["recipe"].(map[string]any)
	require.Equal(t, "Coffee", recipe["title"])
	require.Equal(t, "box", recipe["template"])
	require.NotContains(t, recipe, "logo_data")
	recipe["title"] = "Revised coffee"
	encoded, _ := json.Marshal(recipe)
	resp, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": owner.ID}).Post("/dashboard/media/mockup", strings.NewReader(string(encoded)))
	require.NoError(t, err)
	resp.AssertCreated()
	copied, err := resp.Json()
	require.NoError(t, err)
	copiedAsset := copied["asset"].(map[string]any)
	require.NotEqual(t, asset.ID, copiedAsset["id"])
	t.Cleanup(func() { _ = os.Remove("storage/app/product-media/" + copiedAsset["id"].(string) + ".svg") })
	var original models.ProductMedia
	require.NoError(t, facades.Orm().Query().Where("id = ?", asset.ID).First(&original))
	require.Equal(t, asset.OptionsJSON, original.OptionsJSON)

	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": other.ID}).Get("/dashboard/media/recipe-mockup/recipe")
	require.NoError(t, err)
	resp.AssertNotFound()
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": buyer.ID}).Get("/dashboard/media/recipe-mockup/recipe")
	require.NoError(t, err)
	resp.AssertForbidden()
	resp, err = tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard/media/recipe-photo/recipe")
	require.NoError(t, err)
	resp.AssertNotFound()
}
