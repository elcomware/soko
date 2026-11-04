package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDashboardUsersPaginationStrictAccess(t *testing.T) {
	actor := models.User{ID: "wu-http-admin", Email: "wu-http@test.invalid", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=users&q=wu-http&selected_user=" + actor.ID)
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	require.NotNil(t, props["userPage"])
	require.Equal(t, "available", props["selectedUserStatus"])
	selected := props["selectedUser"].(map[string]any)
	_, secret := selected["password"]
	require.False(t, secret)
	for _, key := range []string{"products", "offers", "orders", "rfqs"} {
		require.Empty(t, props[key])
	}
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=users&role=unsafe")
	require.NoError(t, err)
	response.AssertStatus(400)
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"},{"permission":"users.read","scope":"own"}]`, Version: 1}))
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=users&selected_user=" + actor.ID)
	require.NoError(t, err)
	response.AssertForbidden()
}

func TestDashboardUsersQueryFailureReturnsServerError(t *testing.T) {
	actor := models.User{ID: "wu-http-error", Email: "wu-http-error@test.invalid", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&actor))
	// Authentication can still read the user, but the directory's explicit safe
	// projection must fail rather than claiming an empty successful page.
	_, err := facades.Orm().Query().Exec("ALTER TABLE users RENAME COLUMN country TO wu_unavailable_country")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE users RENAME COLUMN wu_unavailable_country TO country")
		require.NoError(t, e)
	}()
	response, err := (&frameworktesting.TestCase{}).Http(t).WithSession(map[string]any{"user_id": actor.ID}).Get("/dashboard?section=users")
	require.NoError(t, err)
	response.AssertStatus(500)
}
