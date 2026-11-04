package middleware

import (
	"elcomware/app/models"
	"testing"
)

func TestWorkspaceNavigationUsesEffectiveCapabilities(t *testing.T) {
	media := &models.User{ID: "media-only", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "media.create", Scope: "own"}}}
	if !media.CanAccessWorkspace() {
		t.Fatal("media-only staff cannot enter workspace")
	}
	links := WorkspaceNavigation(media)
	if len(links) != 2 || links[1].Href != "/dashboard/media-library" {
		t.Fatal("unrelated navigation exposed", links)
	}
	revoked := &models.User{ID: "revoked", Role: models.RoleAdmin, PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{}}
	if len(WorkspaceNavigation(revoked)) != 0 {
		t.Fatal("role bypassed explicit revocation")
	}
	content := &models.User{ID: "editor", PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "content.update", Scope: "all"}}}
	links = WorkspaceNavigation(content)
	if len(links) != 2 || links[1].Href != "/dashboard/merchandising" {
		t.Fatal(links)
	}
	buyer := &models.User{ID: "buyer", Role: models.RoleBuyer}
	if len(WorkspaceNavigation(buyer)) != 0 {
		t.Fatal("buyer given admin navigation")
	}
}
