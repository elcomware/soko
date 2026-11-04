package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestDashboardPromotionsPaginationScopeAndFailures(t *testing.T) {
	actor := models.User{ID: "wp-http", Email: "wp-http@test.invalid", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard/promotions?per_page=1&selected_campaign=missing")
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	require.NotNil(t, props["campaignPage"])
	require.Equal(t, true, props["targetingApi"])
	require.Empty(t, props["products"])
	require.NotNil(t, props["merchants"])
	require.Equal(t, "unavailable", props["selectedCampaignStatus"])
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard/promotions?status=unsafe")
	require.NoError(t, err)
	response.AssertStatus(400)
	_, err = facades.Orm().Query().Exec("ALTER TABLE promotions RENAME TO wp_unavailable_promotions")
	require.NoError(t, err)
	func() {
		defer func() {
			_, e := facades.Orm().Query().Exec("ALTER TABLE wp_unavailable_promotions RENAME TO promotions")
			require.NoError(t, e)
		}()
		response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard/promotions")
		require.NoError(t, err)
		response.AssertStatus(500)
	}()
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"}]`, Version: 1}))
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard/promotions")
	require.NoError(t, err)
	response.AssertForbidden()
}

func TestPromotionSaveReturnFiltersAreBounded(t *testing.T) {
	actor := models.User{ID: "wp-return", Email: "wp-return@test.invalid", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	response, err := (&frameworktesting.TestCase{}).Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": actor.ID}).Post("/dashboard/promotions", strings.NewReader(`{"return_q":"one & two","return_status":"active","return_sort":"oldest","return_page":"2","return_per_page":"10"}`))
	require.NoError(t, err)
	response.AssertStatus(303).AssertHeader("Location", "/dashboard/promotions?page=2&per_page=10&q=one+%26+two&sort=oldest&status=active")
	response, err = (&frameworktesting.TestCase{}).Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": actor.ID}).Post("/dashboard/promotions", strings.NewReader(`{"return_page":"https://evil.invalid"}`))
	require.NoError(t, err)
	response.AssertStatus(303).AssertHeader("Location", "/dashboard/promotions")
}
