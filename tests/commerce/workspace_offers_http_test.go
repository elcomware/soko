package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDashboardOffersPaginationAndErrors(t *testing.T) {
	actor := models.User{ID: "wof-http", Email: "wof-http@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=offers&selected_offer=missing")
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	require.NotNil(t, props["offerPage"])
	require.Equal(t, "supplier", props["offerVisibility"])
	require.Equal(t, "unavailable", props["selectedOfferStatus"])
	for _, key := range []string{"products", "rfqs", "orders", "users"} {
		require.Empty(t, props[key])
	}
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=offers&status=unsafe")
	require.NoError(t, err)
	response.AssertStatus(400)
	_, err = facades.Orm().Query().Exec("ALTER TABLE offers RENAME TO wof_unavailable_offers")
	require.NoError(t, err)
	func() {
		defer func() {
			_, e := facades.Orm().Query().Exec("ALTER TABLE wof_unavailable_offers RENAME TO offers")
			require.NoError(t, e)
		}()
		response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=offers")
		require.NoError(t, err)
		response.AssertStatus(500)
	}()
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"}]`, Version: 1}))
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=offers")
	require.NoError(t, err)
	response.AssertForbidden()
}
