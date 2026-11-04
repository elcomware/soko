package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"encoding/json"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

// Broad read access must not turn own-scoped mutation grants into global actions.
func TestDashboardRecordCapabilitiesRespectOwnershipAndState(t *testing.T) {
	u := models.User{ID: "record-cap-actor", Email: "record-cap@test.invalid", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&u))
	grants := []models.PermissionGrant{}
	for _, key := range []string{"products.read", "orders.read", "offers.read", "rfqs.read"} {
		grants = append(grants, models.PermissionGrant{Permission: key, Scope: "all"})
	}
	for _, key := range []string{"products.update", "products.delete", "media.create", "orders.cancel", "offers.decide", "rfqs.quote", "rfqs.award", "rfqs.close"} {
		grants = append(grants, models.PermissionGrant{Permission: key, Scope: "own"})
	}
	raw, _ := json.Marshal(grants)
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: u.ID, GrantsJSON: string(raw), Version: 1}))
	for _, p := range []models.Product{{ID: "record-pristine", VendorID: u.ID}, {ID: "record-own", VendorID: u.ID}, {ID: "record-other", VendorID: "other"}, {ID: "record-history", VendorID: u.ID}} {
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	for _, o := range []models.Offer{{ID: "record-offer-own", ProductID: "record-history", Status: models.OfferStatusPending}, {ID: "record-offer-other", ProductID: "record-other", Status: models.OfferStatusPending}, {ID: "record-offer-done", ProductID: "record-own", Status: models.OfferStatusAccepted}} {
		require.NoError(t, facades.Orm().Query().Create(&o))
	}
	for _, o := range []models.Order{{ID: "record-order-own", OrderNumber: "record-order-own", UserID: u.ID, Status: models.OrderStatusPending}, {ID: "record-order-other", OrderNumber: "record-order-other", UserID: "other", Status: models.OrderStatusPending}, {ID: "record-order-done", OrderNumber: "record-order-done", UserID: u.ID, Status: models.OrderStatusCancelled}} {
		require.NoError(t, facades.Orm().Query().Create(&o))
	}
	for _, r := range []models.Rfq{{ID: "record-rfq-own", RfqNumber: "record-rfq-own", BuyerID: u.ID, Status: models.RfqStatusOpen}, {ID: "record-rfq-other", RfqNumber: "record-rfq-other", BuyerID: "other", Status: models.RfqStatusOpen}, {ID: "record-rfq-done", RfqNumber: "record-rfq-done", BuyerID: u.ID, Status: models.RfqStatusClosed}} {
		require.NoError(t, facades.Orm().Query().Create(&r))
	}
	for _, q := range []models.RfqQuote{{ID: "record-quote-pending", RfqID: "record-rfq-own", Status: "pending"}, {ID: "record-quote-done", RfqID: "record-rfq-own", Status: "rejected"}} {
		require.NoError(t, facades.Orm().Query().Create(&q))
	}
	tc := &frameworktesting.TestCase{}
	read := func() map[string]any {
		response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": u.ID}).Get("/dashboard")
		require.NoError(t, err)
		response.AssertOk()
		body, err := response.Json()
		require.NoError(t, err)
		caps, ok := body["props"].(map[string]any)["recordCapabilities"].(map[string]any)
		require.True(t, ok, "dashboard must expose record-specific capabilities")
		return caps
	}
	caps := read()
	for _, row := range []struct {
		kind, id string
		want     map[string]any
	}{
		{"products", "record-pristine", map[string]any{"update": true, "delete": true, "media": true}},
		{"products", "record-own", map[string]any{"update": true, "delete": false, "media": true}},
		{"products", "record-history", map[string]any{"update": true, "delete": false, "media": true}},
		{"products", "record-other", map[string]any{"update": false, "delete": false, "media": false}},
		{"orders", "record-order-own", map[string]any{"cancel": true}},
		{"orders", "record-order-other", map[string]any{"cancel": false}},
		{"orders", "record-order-done", map[string]any{"cancel": false}},
		{"offers", "record-offer-own", map[string]any{"decide": true}},
		{"offers", "record-offer-other", map[string]any{"decide": false}},
		{"offers", "record-offer-done", map[string]any{"decide": false}},
		{"rfqs", "record-rfq-own", map[string]any{"quote": false, "award": true, "close": true}},
		{"rfqs", "record-rfq-other", map[string]any{"quote": true, "award": false, "close": false}},
		{"rfqs", "record-rfq-done", map[string]any{"quote": false, "award": false, "close": false}},
		{"quotes", "record-quote-pending", map[string]any{"accept": true}},
		{"quotes", "record-quote-done", map[string]any{"accept": false}},
	} {
		require.Equal(t, row.want, caps[row.kind].(map[string]any)[row.id], row.kind+"/"+row.id)
	}

	// Global update permission must not imply permission to create media assets.
	raw, _ = json.Marshal([]models.PermissionGrant{{Permission: "products.read", Scope: "all"}, {Permission: "products.update", Scope: "all"}})
	_, err := facades.Orm().Query().Model(&models.AccessProfile{}).Where("user_id = ?", u.ID).Update("grants_json", string(raw))
	require.NoError(t, err)
	caps = read()
	require.Equal(t, map[string]any{"update": true, "delete": false, "media": false}, caps["products"].(map[string]any)["record-other"])
	// Explicit global grants may act on other owners' records; role labels alone never suffice.
	raw, _ = json.Marshal(models.DefaultPermissionGrants(models.RoleAdmin))
	_, err = facades.Orm().Query().Model(&models.AccessProfile{}).Where("user_id = ?", u.ID).Update("grants_json", string(raw))
	require.NoError(t, err)
	caps = read()
	require.Equal(t, map[string]any{"decide": true}, caps["offers"].(map[string]any)["record-offer-other"])
	require.Equal(t, map[string]any{"cancel": true}, caps["orders"].(map[string]any)["record-order-other"])
	require.Equal(t, map[string]any{"quote": true, "award": true, "close": true}, caps["rfqs"].(map[string]any)["record-rfq-other"])
	// Offer decisions must still resolve ownership when products.read is revoked.
	raw, _ = json.Marshal([]models.PermissionGrant{{Permission: "offers.read", Scope: "all"}, {Permission: "offers.decide", Scope: "own"}})
	_, err = facades.Orm().Query().Model(&models.AccessProfile{}).Where("user_id = ?", u.ID).Update("grants_json", string(raw))
	require.NoError(t, err)
	caps = read()
	require.Empty(t, caps["products"])
	require.Equal(t, map[string]any{"decide": true}, caps["offers"].(map[string]any)["record-offer-own"])
	require.Equal(t, map[string]any{"decide": false}, caps["offers"].(map[string]any)["record-offer-other"])
}
