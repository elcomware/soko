package models

import "testing"

func TestGranularPermissions(t *testing.T) {
	u := &User{ID: "m", Role: RoleAdmin, PermissionsLoaded: true, PermissionGrants: []PermissionGrant{{"products.read", "own"}}}
	if !u.Can("products.read", "m") || u.Can("products.read", "other") || u.Can("products.read", "") || u.Can("products.update", "m") || u.Can("unknown", "m") {
		t.Fatal("scope boundary failed")
	}
	u.PermissionGrants = nil
	if u.Can("permissions.manage", "") {
		t.Fatal("role bypassed explicit revocation")
	}
	if (&User{ID: "x", Role: "unknown"}).HasPermission("products.read") {
		t.Fatal("unknown role granted access")
	}
}
func TestPermissionTemplates(t *testing.T) {
	for _, r := range []UserRole{RoleAdmin, RoleMerchant, RoleBuyer} {
		for _, g := range DefaultPermissionGrants(r) {
			if !ValidPermission(g.Permission) || (g.Scope != "own" && g.Scope != "all") {
				t.Fatal("bad default")
			}
		}
	}
	u := &User{ID: "editor", PermissionsLoaded: true, PermissionGrants: []PermissionGrant{{"content.update", "all"}}}
	if !u.CanAccessWorkspace() || u.CanManageCatalog() {
		t.Fatal("content editor should access workspace without catalog write")
	}
}
