package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
	"elcomware/app/services"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestExplicitProfileRevokesRoleAndScopesWorkspace(t *testing.T) {
	admin := &models.User{ID: "access-admin", Name: "Admin", Email: "access-admin@test.invalid", Role: models.RoleAdmin}
	target := &models.User{ID: "access-editor", Name: "Editor", Email: "access-editor@test.invalid", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(admin))
	require.NoError(t, facades.Orm().Query().Create(target))
	require.NoError(t, services.SaveAccessProfile(admin, target.ID, []models.PermissionGrant{{Permission: "content.update", Scope: "all"}}, 0))
	queries.HydrateAccess(target)
	require.True(t, target.CanAccessWorkspace())
	require.False(t, target.CanManageCatalog())
	data, err := services.LoadWorkspace(target, false)
	require.NoError(t, err)
	require.Empty(t, data.Products)
	require.Empty(t, data.Users)
	require.Empty(t, data.Orders)
	require.Error(t, services.SaveAccessProfile(admin, target.ID, nil, 0))
	require.NoError(t, services.SaveAccessProfile(admin, target.ID, []models.PermissionGrant{}, 1))
	queries.HydrateAccess(target)
	require.False(t, target.CanAccessWorkspace())
	require.Error(t, services.SaveAccessProfile(admin, admin.ID, nil, 0))
	require.Error(t, services.SaveAccessProfile(target, admin.ID, nil, 0))
	count, err := facades.Orm().Query().Model(&models.AccessAudit{}).Where("target_id = ?", target.ID).Count()
	require.NoError(t, err)
	require.EqualValues(t, 2, count)
}
func TestBrandAndPaymentVersionedSettings(t *testing.T) {
	admin := &models.User{ID: "settings-admin", Role: models.RoleAdmin}
	merchant := &models.User{ID: "settings-merchant", Role: models.RoleMerchant}
	b := queries.DefaultBrandSettings()
	b.DisplayName = "Soko studio"
	b.CardFooterPrice = false
	b.CardFooterRating = true
	b.CardFooterMode = "rotate"
	b.CardFooterIntervalSeconds = 12
	require.Error(t, services.SaveBrandSettings(merchant, b))
	require.NoError(t, services.SaveBrandSettings(admin, b))
	require.Error(t, services.SaveBrandSettings(admin, b))
	got := queries.LoadBrandSettings()
	require.Equal(t, "Soko studio", got.DisplayName)
	require.False(t, got.CardFooterPrice)
	require.True(t, got.CardFooterRating)
	require.True(t, got.CardFooterAvailability)
	require.Equal(t, "rotate", got.CardFooterMode)
	require.Equal(t, 12, got.CardFooterIntervalSeconds)
	require.Equal(t, 1, got.Version)
	p := services.PaymentPreference{Provider: "stripe", Methods: []string{"card"}, Country: "US", Currency: "USD", Routes: []services.PaymentRoutePreference{{Country: "KE", Area: "KE-30", Currency: "KES", Method: "card", Providers: []string{"stripe", "flutterwave"}}}}
	require.Error(t, services.SavePaymentPreference(merchant, "someone-else", p))
	require.NoError(t, services.SavePaymentPreference(merchant, merchant.ID, p))
	stored := services.LoadPaymentPreference(merchant.ID)
	require.False(t, stored.Enabled)
	require.Equal(t, p.Routes, stored.Routes)
	require.Equal(t, "requires_provider_onboarding", stored.Status)
	require.Error(t, services.SavePaymentPreference(merchant, merchant.ID, p))
	p.Version = 1
	p.Enabled = true
	require.Error(t, services.SavePaymentPreference(merchant, merchant.ID, p))
}

func TestHTTPRevokedAdminCannotReadOrWriteSettings(t *testing.T) {
	u := models.User{ID: "access-http-revoked", Name: "Former admin", Email: "access-http-revoked@test.invalid", Role: models.RoleAdmin}
	require.NoError(t, facades.Orm().Query().Create(&u))
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: u.ID, GrantsJSON: "[]", Version: 1}))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithSession(map[string]any{"user_id": u.ID, "user_role": "admin"}).Get("/dashboard/settings")
	require.NoError(t, err)
	response.AssertForbidden()
	response, err = tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"}).WithSession(map[string]any{"user_id": u.ID, "user_role": "admin"}).Post("/dashboard/content", strings.NewReader(`{}`))
	require.NoError(t, err)
	response.AssertForbidden()
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": u.ID}).Get("/dashboard")
	require.NoError(t, err)
	response.AssertFound().AssertHeader("Location", "/account")
}

func TestDelegationCeilingAndAdministratorRetention(t *testing.T) {
	a := models.User{ID: "access-admin-a", Email: "access-admin-a@test.invalid", Role: models.RoleAdmin}
	b := models.User{ID: "access-admin-b", Email: "access-admin-b@test.invalid", Role: models.RoleAdmin}
	target := models.User{ID: "access-target", Email: "access-target@test.invalid", Role: models.RoleBuyer}
	for _, u := range []*models.User{&a, &b, &target} {
		require.NoError(t, facades.Orm().Query().Create(u))
	}
	delegate := &models.User{ID: "delegate", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "permissions.manage", Scope: "all"}}}
	require.Error(t, services.SaveAccessProfile(delegate, target.ID, []models.PermissionGrant{{Permission: "products.update", Scope: "all"}}, 0))
	require.NoError(t, services.SaveAccessProfile(&a, b.ID, []models.PermissionGrant{}, 0))
	require.Error(t, services.SaveAccessProfile(&b, a.ID, []models.PermissionGrant{}, 0))
	queries.HydrateAccess(&a)
	require.True(t, a.Can("permissions.manage", ""))
	queries.HydrateAccess(&b)
	require.False(t, b.Can("permissions.manage", ""))
}
func TestExplicitRevocationHidesPersonalAccountResources(t *testing.T) {
	u := &models.User{ID: "account-revoked", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{}}
	require.NoError(t, facades.Orm().Query().Create(&models.Order{ID: "account-revoked-order", UserID: u.ID, Status: models.OrderStatusPending}))
	d, err := services.LoadWorkspace(u, true)
	require.NoError(t, err)
	require.Empty(t, d.Orders)
	require.Empty(t, d.Offers)
	require.Empty(t, d.Rfqs)
}
func TestAccessStorageFailureFailsClosed(t *testing.T) {
	_, err := facades.Orm().Query().Exec("ALTER TABLE user_access_profiles RENAME TO unavailable_access_profiles")
	require.NoError(t, err)
	defer func() {
		_, err := facades.Orm().Query().Exec("ALTER TABLE unavailable_access_profiles RENAME TO user_access_profiles")
		require.NoError(t, err)
	}()
	admin := &models.User{ID: "storage-failure-admin", Role: models.RoleAdmin}
	queries.HydrateAccess(admin)
	require.True(t, admin.PermissionsLoaded)
	require.False(t, admin.Can("permissions.manage", ""))
	require.False(t, admin.CanAccessWorkspace())
}
func TestDefaultTemplatesValidateWithoutDuplicates(t *testing.T) {
	for _, role := range []models.UserRole{models.RoleBuyer, models.RoleMerchant, models.RoleAdmin} {
		require.NoError(t, services.ValidateGrants(models.DefaultPermissionGrants(role)))
	}
}

func TestMerchantPaymentRoutesAreScopedAndVersioned(t *testing.T) {
	owner := models.User{ID: "route-owner", Email: "route-owner@test.invalid", Role: models.RoleMerchant}
	other := models.User{ID: "route-other", Email: "route-other@test.invalid", Role: models.RoleMerchant}
	admin := models.User{ID: "route-admin", Email: "route-admin@test.invalid", Role: models.RoleAdmin}
	for _, u := range []*models.User{&owner, &other, &admin} {
		require.NoError(t, facades.Orm().Query().Create(u))
	}
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithSession(map[string]any{"user_id": owner.ID}).Get("/dashboard/settings?merchant_id=" + other.ID)
	require.NoError(t, err)
	response.AssertForbidden()
	payload := `{"provider":"none","methods":[],"country":"KE","currency":"KES","routes":[{"country":"KE","area":"","currency":"KES","method":"mobile_money","providers":["flutterwave"]}]}`
	response, err = tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"}).WithSession(map[string]any{"user_id": owner.ID}).Post("/dashboard/settings/payments?merchant_id="+other.ID, strings.NewReader(payload))
	require.NoError(t, err)
	response.AssertForbidden()
	response, err = tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/settings/payments?merchant_id="+other.ID, strings.NewReader(payload))
	require.NoError(t, err)
	response.AssertHeader("Location", "/dashboard/settings?merchant_id="+other.ID)
	saved := services.LoadPaymentPreference(other.ID)
	require.Len(t, saved.Routes, 1)
	require.Equal(t, 1, saved.Version)
	require.False(t, saved.Enabled)
	require.Equal(t, "requires_provider_onboarding", saved.Status)
	stale := saved
	stale.Version = 0
	stale.Routes = nil
	require.Error(t, services.SavePaymentPreference(&admin, other.ID, stale))
	require.Len(t, services.LoadPaymentPreference(other.ID).Routes, 1)
	saved.Enabled = true
	require.Error(t, services.SavePaymentPreference(&admin, other.ID, saved))
}
