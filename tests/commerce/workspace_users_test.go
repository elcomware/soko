package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWorkspaceUsersStrictScopePagingAndSelection(t *testing.T) {
	admin := &models.User{ID: "wu-admin", Role: models.RoleAdmin}
	stamp := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, u := range []models.User{{ID: "wu-a", Email: "wu-a@test.invalid", Name: "wu-scope ÉCOLE", Role: models.RoleBuyer, Verified: true, Password: "not-for-response", CreatedAt: stamp}, {ID: "wu-b", Email: "wu-b@test.invalid", Name: "wu-scope Buyer", Role: models.RoleMerchant, CreatedAt: stamp}} {
		require.NoError(t, facades.Orm().Query().Create(&u))
	}
	f, err := services.ParseWorkspaceUserFilters(map[string]string{"q": "wu-scope", "per_page": "1", "selected_user": "wu-a"})
	require.NoError(t, err)
	page, err := services.LoadWorkspaceUsers(admin, f)
	require.NoError(t, err)
	require.EqualValues(t, 2, page.UserPage.FilteredTotal)
	require.Equal(t, "wu-b", page.Users[0].ID)
	require.Equal(t, "available", page.SelectedUserStatus)
	require.Empty(t, page.SelectedUser.Password)
	require.Empty(t, page.Users[0].Password)
	f.Query = "école"
	f.Page = 999
	page, err = services.LoadWorkspaceUsers(admin, f)
	require.NoError(t, err)
	require.Equal(t, 1, page.UserPage.Page)
	require.EqualValues(t, 1, page.UserPage.FilteredTotal)
	f.Query = "wu-scope"
	f.Role = "merchant"
	page, err = services.LoadWorkspaceUsers(admin, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.UserPage.FilteredTotal)
	require.Equal(t, "wu-b", page.Users[0].ID)
	f.Role = ""
	f.Verification = "verified"
	page, err = services.LoadWorkspaceUsers(admin, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.UserPage.FilteredTotal)
	require.Equal(t, "wu-a", page.Users[0].ID)
	f.Verification = ""
	f.Sort = "oldest"
	f.Page = 1
	f.SelectedUser = "wu-missing"
	page, err = services.LoadWorkspaceUsers(admin, f)
	require.NoError(t, err)
	require.Equal(t, "wu-a", page.Users[0].ID)
	require.Nil(t, page.SelectedUser)
	require.Equal(t, "unavailable", page.SelectedUserStatus)
	own := &models.User{ID: "wu-a", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "workspace.access", Scope: "own"}, {Permission: "users.read", Scope: "own"}}}
	page, err = services.LoadWorkspaceUsers(own, f)
	require.Error(t, err)
	require.Nil(t, page)
}

func TestWorkspaceUsersUnicodeVerificationAndErrors(t *testing.T) {
	actor := &models.User{ID: "wu-reader", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "users.read", Scope: "all"}}}
	for i, name := range []string{"wu-u ΟΣ", "wu-u ПРИВЕТ", "wu-u İSTANBUL", "wu-u 100%_"} {
		id := fmt.Sprintf("wu-u-%d", i)
		require.NoError(t, facades.Orm().Query().Create(&models.User{ID: id, Email: id + "@test.invalid", Name: name, Role: models.RoleBuyer, Verified: i%2 == 0}))
	}
	for _, q := range []string{"wu-u ος", "wu-u привет", "wu-u i\u0307stanbul", "wu-u 100%_"} {
		f, err := services.ParseWorkspaceUserFilters(map[string]string{"q": q})
		require.NoError(t, err)
		page, err := services.LoadWorkspaceUsers(actor, f)
		require.NoError(t, err)
		require.EqualValues(t, 1, page.UserPage.FilteredTotal, q)
	}
	for _, q := range []string{"VERIFIED", "unverified"} {
		f, err := services.ParseWorkspaceUserFilters(map[string]string{"q": q})
		require.NoError(t, err)
		page, err := services.LoadWorkspaceUsers(actor, f)
		require.NoError(t, err)
		for _, u := range page.Users {
			require.Equal(t, q == "VERIFIED", u.Verified)
		}
	}
	f, err := services.ParseWorkspaceUserFilters(map[string]string{"q": "verified", "verification": "unverified"})
	require.NoError(t, err)
	page, err := services.LoadWorkspaceUsers(actor, f)
	require.NoError(t, err)
	require.Zero(t, page.UserPage.FilteredTotal)
	f, err = services.ParseWorkspaceUserFilters(map[string]string{"selected_user": "wu-u-0"})
	require.NoError(t, err)
	_, err = facades.Orm().Query().Model(&models.User{}).Where("id = ?", "wu-u-0").Delete()
	require.NoError(t, err)
	page, err = services.LoadWorkspaceUsers(actor, f)
	require.NoError(t, err)
	require.Nil(t, page.SelectedUser)
	require.Equal(t, "unavailable", page.SelectedUserStatus)
	_, err = facades.Orm().Query().Exec("ALTER TABLE users RENAME TO wu_unavailable_users")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE wu_unavailable_users RENAME TO users")
		require.NoError(t, e)
	}()
	page, err = services.LoadWorkspaceUsers(actor, f)
	require.Error(t, err)
	require.Nil(t, page)
}
func TestWorkspaceUserFilterValidation(t *testing.T) {
	for _, input := range []map[string]string{{"role": "root"}, {"verification": "yes"}, {"sort": "email"}, {"page": "0"}, {"per_page": "101"}, {"q": "bad\nquery"}} {
		_, err := services.ParseWorkspaceUserFilters(input)
		require.Error(t, err)
	}
	f, err := services.ParseWorkspaceUserFilters(nil)
	require.NoError(t, err)
	_, err = services.LoadWorkspaceUsers(nil, f)
	require.Error(t, err)
}
