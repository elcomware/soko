package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Dashboard actions must preserve context without widening existing domain permissions.
func TestDashboardOrderAndRfqActionsKeepContextAndPermissions(t *testing.T) {
	admin := models.User{ID: "dash-action-admin", Email: "dash-action-admin@test.invalid", Role: models.RoleAdmin}
	seller := models.User{ID: "dash-action-seller", Email: "dash-action-seller@test.invalid", Role: models.RoleMerchant}
	buyer := models.User{ID: "dash-action-buyer", Email: "dash-action-buyer@test.invalid", Role: models.RoleBuyer}
	revoked := models.User{ID: "dash-action-revoked", Email: "dash-action-revoked@test.invalid", Role: models.RoleAdmin}
	for _, u := range []*models.User{&admin, &seller, &buyer, &revoked} {
		require.NoError(t, facades.Orm().Query().Create(u))
	}
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: revoked.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"},{"permission":"orders.read","scope":"all"},{"permission":"rfqs.read","scope":"all"}]`, Version: 1}))
	p := models.Product{ID: "dash-action-product", VendorID: seller.ID, Stock: 3, Price: 10, MOQ: 1, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&p))
	order := models.Order{ID: "dash-action-order", OrderNumber: "dash-action-order", UserID: buyer.ID, Status: models.OrderStatusPending}
	require.NoError(t, facades.Orm().Query().Create(&order))
	require.NoError(t, facades.Orm().Query().Create(&models.OrderItem{ID: "dash-action-item", OrderID: order.ID, ProductID: p.ID, Quantity: 2, StockReserved: true}))
	for _, id := range []string{"dash-action-close", "dash-action-accept"} {
		require.NoError(t, facades.Orm().Query().Create(&models.Rfq{ID: id, RfqNumber: id, BuyerID: buyer.ID, Status: models.RfqStatusOpen}))
	}
	require.NoError(t, facades.Orm().Query().Create(&models.RfqQuote{ID: "dash-action-quote", RfqID: "dash-action-accept", VendorID: seller.ID, Status: "pending", UpdatedAt: time.Now()}))
	tc := &frameworktesting.TestCase{}
	actions := []struct{ path, dest string }{
		{"/dashboard/orders/" + order.ID + "/cancel", "/dashboard?section=orders"},
		{"/dashboard/rfqs/dash-action-close/close", "/dashboard?section=rfqs"},
		{"/dashboard/rfqs/dash-action-accept/quotes/dash-action-quote/accept", "/dashboard?section=rfqs"},
	}
	for _, a := range actions {
		response, err := tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").Post(a.path, strings.NewReader(""))
		require.NoError(t, err)
		response.AssertFound().AssertHeader("Location", "/login")
		// Buyer owns these resources but has no operational workspace access.
		response, err = tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": buyer.ID}).Post(a.path, strings.NewReader(""))
		require.NoError(t, err)
		response.AssertForbidden()
		for _, id := range []string{seller.ID, revoked.ID} {
			response, err = tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": id}).Post(a.path, strings.NewReader(""))
			require.NoError(t, err)
			if id == seller.ID && a.path == actions[0].path {
				response.AssertForbidden()
			} else {
				response.AssertStatus(303).AssertHeader("Location", a.dest)
			}
		}
	}
	var before models.Order
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.ID).First(&before))
	require.Equal(t, models.OrderStatusPending, before.Status)
	for _, id := range []string{"dash-action-close", "dash-action-accept"} {
		var rfq models.Rfq
		require.NoError(t, facades.Orm().Query().Where("id = ?", id).First(&rfq))
		require.Equal(t, models.RfqStatusOpen, rfq.Status)
	}

	// A partial supplier view may not cancel unseen order lines, even with cancel/all.
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: seller.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"},{"permission":"orders.read","scope":"own"},{"permission":"orders.cancel","scope":"all"}]`, Version: 1}))
	partial, err := tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": seller.ID}).Post(actions[0].path, strings.NewReader(""))
	require.NoError(t, err)
	partial.AssertForbidden()
	// Switching to the buyer endpoint must not bypass whole-order visibility.
	bypass, err := tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": seller.ID}).Post("/account/orders/"+order.ID+"/cancel", strings.NewReader(""))
	require.NoError(t, err)
	bypass.AssertStatus(303).AssertHeader("Location", "/account?section=orders")
	var unchanged models.Order
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.ID).First(&unchanged))
	require.Equal(t, models.OrderStatusPending, unchanged.Status)
	view, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": seller.ID}).Get("/dashboard?section=orders")
	require.NoError(t, err)
	view.AssertOk()
	page, err := view.Json()
	require.NoError(t, err)
	caps := page["props"].(map[string]any)["recordCapabilities"].(map[string]any)["orders"].(map[string]any)
	require.Equal(t, false, caps[order.ID].(map[string]any)["cancel"])
	// A quote belonging to another request must never be accepted by changing the URL.
	response, err := tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/rfqs/dash-action-close/quotes/dash-action-quote/accept", strings.NewReader(""))
	require.NoError(t, err)
	response.AssertStatus(303).AssertHeader("Location", "/dashboard?section=rfqs")
	var quote models.RfqQuote
	require.NoError(t, facades.Orm().Query().Where("id = ?", "dash-action-quote").First(&quote))
	require.Equal(t, "pending", quote.Status)
	for _, a := range actions {
		body := url.Values{"quote_version": {quote.UpdatedAt.Format(time.RFC3339Nano)}}.Encode()
		response, err := tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/x-www-form-urlencoded"}).WithSession(map[string]any{"user_id": admin.ID}).Post(a.path, strings.NewReader(body))
		require.NoError(t, err)
		response.AssertStatus(303).AssertHeader("Location", a.dest)
	}
	var after models.Order
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.ID).First(&after))
	require.Equal(t, models.OrderStatusCancelled, after.Status)
	var closed, awarded models.Rfq
	require.NoError(t, facades.Orm().Query().Where("id = ?", "dash-action-close").First(&closed))
	require.NoError(t, facades.Orm().Query().Where("id = ?", "dash-action-accept").First(&awarded))
	require.Equal(t, models.RfqStatusClosed, closed.Status)
	require.Equal(t, models.RfqStatusAwarded, awarded.Status)
	require.NoError(t, facades.Orm().Query().Where("id = ?", "dash-action-quote").First(&quote))
	require.Equal(t, "accepted", quote.Status)
	// A retry cannot release inventory twice, and cannot change the return context.
	response, err = tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": admin.ID}).Post(actions[0].path+"?return_to=https://untrusted.invalid", strings.NewReader(""))
	require.NoError(t, err)
	response.AssertStatus(303).AssertHeader("Location", actions[0].dest)
	var stock models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", p.ID).First(&stock))
	require.Equal(t, 5, stock.Stock)
}
