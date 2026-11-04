package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDashboardRfqsPaginationAndErrors(t *testing.T) {
	actor := models.User{ID: "wr-http", Email: "wr-http@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=rfqs&selected_rfq=missing")
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	require.NotNil(t, props["rfqPage"])
	require.Equal(t, "supplier", props["rfqVisibility"])
	require.Equal(t, "unavailable", props["selectedRfqStatus"])
	for _, key := range []string{"products", "offers", "orders", "users"} {
		require.Empty(t, props[key])
	}
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=rfqs&status=unsafe")
	require.NoError(t, err)
	response.AssertStatus(400)
	_, err = facades.Orm().Query().Exec("ALTER TABLE rfqs RENAME TO wr_unavailable_rfqs")
	require.NoError(t, err)
	func() {
		defer func() {
			_, e := facades.Orm().Query().Exec("ALTER TABLE wr_unavailable_rfqs RENAME TO rfqs")
			require.NoError(t, e)
		}()
		response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=rfqs")
		require.NoError(t, err)
		response.AssertStatus(500)
	}()
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"}]`, Version: 1}))
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=rfqs")
	require.NoError(t, err)
	response.AssertForbidden()
}
