package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"encoding/json"
)

// HydrateAccess fails closed on storage/decoding errors. An absent profile uses the
// onboarding template; explicit empty profiles revoke every permission.
func HydrateAccess(u *models.User) {
	if u == nil {
		return
	}
	u.PermissionsLoaded = true
	u.PermissionGrants = nil

	var p models.AccessProfile
	if err := facades.Orm().Query().Where("user_id = ?", u.ID).First(&p); err != nil {
		return
	}
	if p.UserID == "" {
		u.PermissionGrants = models.DefaultPermissionGrants(u.Role)
		return
	}
	var grants []models.PermissionGrant
	if json.Unmarshal([]byte(p.GrantsJSON), &grants) == nil {
		u.PermissionGrants = grants
	}
}
