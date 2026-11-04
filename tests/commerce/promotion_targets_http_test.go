package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestPromotionTargetsEndpointReadOnlyAndScoped(t *testing.T) {
	actor := models.User{ID: "pt-http", Email: "pt-http@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"promotions.manage","scope":"own"}]`, Version: 1}))
	tc := &frameworktesting.TestCase{}
	post := func(body string, intended int) {
		response, err := tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": actor.ID}).Post("/dashboard/promotions/targets", strings.NewReader(body))
		require.NoError(t, err)
		response.AssertStatus(intended)
	}
	before, err := facades.Orm().Query().Model(&models.Promotion{}).Count()
	require.NoError(t, err)
	post(`{"vendor_id":"pt-http","scope":"{}","per_page":"1"}`, 200)
	post(`{"vendor_id":"foreign","scope":"{}"}`, 403)
	post(`{"vendor_id":"","scope":"{}"}`, 403)
	post(`{"vendor_id":"pt-http","scope":"not json"}`, 400)
	post(`{"vendor_id":"pt-http","scope":"{}","per_page":"101"}`, 400)
	after, err := facades.Orm().Query().Model(&models.Promotion{}).Count()
	require.NoError(t, err)
	require.Equal(t, before, after)
}
