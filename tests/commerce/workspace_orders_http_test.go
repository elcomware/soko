package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDashboardOrdersPaginationAndErrors(t *testing.T) {
	actor := models.User{ID: "wo-http", Email: "wo-http@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=orders&selected_order=missing")
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	require.NotNil(t, props["orderPage"])
	require.Equal(t, "supplier", props["orderVisibility"])
	require.Equal(t, "unavailable", props["selectedOrderStatus"])
	for _, key := range []string{"products", "offers", "rfqs", "users"} {
		require.Empty(t, props[key])
	}
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=orders&status=unsafe")
	require.NoError(t, err)
	response.AssertStatus(400)
	_, err = facades.Orm().Query().Exec("ALTER TABLE orders RENAME TO wo_unavailable_orders")
	require.NoError(t, err)
	func() {
		defer func() {
			_, e := facades.Orm().Query().Exec("ALTER TABLE wo_unavailable_orders RENAME TO orders")
			require.NoError(t, e)
		}()
		response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=orders")
		require.NoError(t, err)
		response.AssertStatus(500)
	}()
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"}]`, Version: 1}))
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=orders")
	require.NoError(t, err)
	response.AssertForbidden()
}
