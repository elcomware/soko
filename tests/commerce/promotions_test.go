package commerce_test

import (
	"elcomware/app/actions"
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"encoding/json"
	"github.com/google/uuid"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestPromotionPermissionsScheduleAndHistoricalOrderPrice(t *testing.T) {
	now := time.Now().UTC()
	admin := &models.User{ID: "promotion-admin", Role: models.RoleAdmin}
	buyer := &models.User{ID: "promotion-buyer", Role: models.RoleBuyer}
	campaign := models.Promotion{Name: "Winter clothing", Enabled: true, Kind: "percentage", Amount: 20, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour), Scope: models.PromotionScope{Categories: []string{"Promo clothing"}, ExcludeProductIDs: []string{"promotion-excluded"}}}
	_, err := services.SavePromotion(campaign, buyer, "")
	require.Error(t, err)
	saved, err := services.SavePromotion(campaign, admin, "")
	require.NoError(t, err)
	p := models.Product{ID: "promotion-shirt", Title: "Shirt", Category: "Promo clothing", Price: 500, MOQ: 1, Stock: 10, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&p))
	order, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{p.ID: 2}, buyer)
	require.NoError(t, err)
	require.Equal(t, 800.0, order.TotalAmount)
	require.Equal(t, 500.0, order.Items[0].OriginalUnitPrice)
	require.Equal(t, 100.0, order.Items[0].DiscountPerUnit)
	require.Equal(t, saved.ID, order.Items[0].PromotionID)
	saved.Enabled = false
	_, err = services.SavePromotion(*saved, admin, "stale")
	require.Error(t, err)
	_, err = services.SavePromotion(*saved, admin, saved.UpdatedAt.Format(time.RFC3339Nano))
	require.NoError(t, err)
	next, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{p.ID: 1}, buyer)
	require.NoError(t, err)
	require.Equal(t, 500.0, next.TotalAmount)
	var historical models.OrderItem
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.Items[0].ID).First(&historical))
	require.Equal(t, 400.0, historical.UnitPrice)
	require.Equal(t, 100.0, historical.DiscountPerUnit)
}
func TestPromotionOwnPermissionCannotAffectOtherMerchant(t *testing.T) {
	merchant := &models.User{ID: "campaign-owner", Role: models.RoleMerchant, PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "promotions.manage", Scope: "own"}}}
	now := time.Now()
	c := models.Promotion{Name: "Own sale", VendorID: merchant.ID, Enabled: true, Kind: "fixed", Amount: 10, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour), Scope: models.PromotionScope{All: true}}
	saved, err := services.SavePromotion(c, merchant, "")
	require.NoError(t, err)
	require.False(t, saved.Applies(models.Product{ID: "other", VendorID: "other"}, now))
	c.VendorID = ""
	_, err = services.SavePromotion(c, merchant, "")
	require.Error(t, err)
	c.VendorID = "other"
	_, err = services.SavePromotion(c, merchant, "")
	require.Error(t, err)
}

func TestHTTPPromotionsRejectUnauthorizedAndPersistExplicitScope(t *testing.T) {
	admin := models.User{ID: "http-promo-admin", Name: "Admin", Email: "promo-admin@example.com", Role: models.RoleAdmin}
	buyer := models.User{ID: "http-promo-buyer", Name: "Buyer", Email: "promo-buyer@example.com", Role: models.RoleBuyer}
	require.NoError(t, facades.Orm().Query().Create(&admin))
	require.NoError(t, facades.Orm().Query().Create(&buyer))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithSession(map[string]any{"user_id": buyer.ID}).Get("/dashboard/promotions")
	require.NoError(t, err)
	response.AssertForbidden()
	response, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": buyer.ID}).Post("/dashboard/promotions", strings.NewReader(`{}`))
	require.NoError(t, err)
	response.AssertForbidden()
	values := map[string]string{"name": "HTTP Holiday", "kind": "percentage", "amount": "15", "enabled": "true", "starts_at": "2026-12-01T00:00:00Z", "ends_at": "2027-01-01T00:00:00Z", "scope": `{"categories":["Clothing"],"exclude_sub_categories":[{"category":"Clothing","sub_category":"Clearance"}]}`}
	body, _ := json.Marshal(values)
	response, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/promotions", strings.NewReader(string(body)))
	require.NoError(t, err)
	response.AssertStatus(303)
	var campaign models.Promotion
	require.NoError(t, facades.Orm().Query().Where("name = ?", "HTTP Holiday").First(&campaign))
	require.NotEmpty(t, campaign.ID)
	campaign.Hydrate()
	require.Equal(t, []string{"Clothing"}, campaign.Scope.Categories)
	require.True(t, campaign.Enabled)
}

func TestHTTPProductQuoteUsesQuantityTierAndLivePromotions(t *testing.T) {
	now := time.Now()
	p := models.Product{ID: "quote-tier-product", Title: "Quote", Price: 100, WholesalePrice: 80, Format: models.FormatWholesale, MOQ: 1, Stock: 100, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&p))
	c := models.Promotion{Name: "Quote sale", Enabled: true, Kind: "percentage", Amount: 20, StartsAt: now.Add(-time.Hour), EndsAt: now.Add(time.Hour), Scope: models.PromotionScope{ProductIDs: []string{p.ID}}}
	saved, err := services.SavePromotion(c, &models.User{ID: "quote-admin", Role: models.RoleAdmin}, "")
	require.NoError(t, err)
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).Get("/products/" + p.ID + "/price?quantity=10")
	require.NoError(t, err)
	response.AssertOk().AssertHeader("Cache-Control", "no-store")
	body, err := response.Json()
	require.NoError(t, err)
	pricing := body["pricing"].(map[string]any)
	require.Equal(t, 80.0, pricing["original_unit_price"])
	require.Equal(t, 64.0, pricing["unit_price"])
	require.Equal(t, saved.ID, pricing["promotion_id"])
	saved.EndsAt = now.Add(-time.Second)
	saved.StartsAt = now.Add(-time.Hour)
	require.NoError(t, facades.Orm().Query().Save(saved))
	response, err = tc.Http(t).Get("/products/" + p.ID + "/price?quantity=10")
	require.NoError(t, err)
	body, err = response.Json()
	require.NoError(t, err)
	require.Equal(t, 80.0, body["pricing"].(map[string]any)["unit_price"])
	for _, qty := range []string{"0", "1.5", "101", "1000001"} {
		response, err = tc.Http(t).Get("/products/" + p.ID + "/price?quantity=" + qty)
		require.NoError(t, err)
		response.AssertStatus(422)
	}
	response, err = tc.Http(t).Get("/products/missing-quote/price?quantity=1")
	require.NoError(t, err)
	response.AssertNotFound()
}

func TestCheckoutTokenRetryNeverReservesStockTwice(t *testing.T) {
	p := models.Product{ID: "idempotent-product", Title: "Idempotent", Price: 10, MOQ: 1, Stock: 10, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&p))
	actor := &models.User{ID: "idempotent-buyer", Role: models.RoleBuyer}
	action := actions.NewExecuteCheckoutAction()
	token := uuid.NewString()
	first, err := action.ExecuteCartWithToken(map[string]int{p.ID: 2}, actor, token)
	require.NoError(t, err)
	// Even a response lost after cart clearing resolves to the original order.
	retry, err := action.ExecuteCartWithToken(map[string]int{}, actor, token)
	require.NoError(t, err)
	require.Equal(t, first.ID, retry.ID)
	var stored models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", p.ID).First(&stored))
	require.Equal(t, 8, stored.Stock)
	_, err = action.ExecuteCartWithToken(map[string]int{p.ID: 1}, actor, "invalid")
	require.Error(t, err)
	// The token is buyer-bound, so it cannot reveal another buyer's order.
	other, err := action.ExecuteCartWithToken(map[string]int{p.ID: 1}, &models.User{ID: "other-idempotent-buyer", Role: models.RoleBuyer}, token)
	require.NoError(t, err)
	require.NotEqual(t, first.ID, other.ID)
}

func TestHTTPPromotionRejectsInvalidMerchantAndAcceptsCatalogMerchant(t *testing.T) {
	admin := models.User{ID: "campaign-picker-admin", Name: "Admin", Email: "campaign-picker-admin@example.com", Role: models.RoleAdmin}
	merchant := models.User{ID: "campaign-picker-merchant", Name: "Merchant", Email: "campaign-picker-merchant@example.com", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&admin))
	require.NoError(t, facades.Orm().Query().Create(&merchant))
	tc := &frameworktesting.TestCase{}
	for _, vendor := range []string{"does-not-exist", merchant.ID} {
		values := map[string]string{"name": "Merchant picker " + vendor, "vendor_id": vendor, "kind": "percentage", "amount": "10", "enabled": "false", "starts_at": "2026-12-01T00:00:00Z", "ends_at": "2027-01-01T00:00:00Z", "scope": `{"all":true}`}
		body, _ := json.Marshal(values)
		response, err := tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/promotions", strings.NewReader(string(body)))
		require.NoError(t, err)
		response.AssertStatus(303)
		var campaign models.Promotion
		require.NoError(t, facades.Orm().Query().Where("name = ?", values["name"]).First(&campaign))
		if vendor == merchant.ID {
			require.NotEmpty(t, campaign.ID)
			require.Equal(t, merchant.ID, campaign.VendorID)
		} else {
			require.Empty(t, campaign.ID)
		}
	}
}
