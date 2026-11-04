package models

// PermissionGrant gives an action either over the actor's own resources or all resources.
type PermissionGrant struct {
	Permission string `json:"permission"`
	Scope      string `json:"scope"`
}
type AccessProfile struct {
	UserID     string `gorm:"primaryKey;type:varchar(64)" json:"user_id"`
	GrantsJSON string `gorm:"type:text" json:"-"`
	Version    int    `json:"version"`
}

func (AccessProfile) TableName() string { return "user_access_profiles" }

type AccessAudit struct {
	ID         string `gorm:"primaryKey;type:varchar(64)"`
	ActorID    string
	TargetID   string
	BeforeJSON string
	AfterJSON  string
	CreatedAt  int64
}

func (AccessAudit) TableName() string { return "access_audits" }

var PermissionKeys = []string{"workspace.access", "products.read", "products.create", "products.update", "products.delete", "orders.read", "orders.cancel", "offers.read", "offers.decide", "rfqs.read", "rfqs.quote", "rfqs.award", "rfqs.close", "users.read", "content.update", "brand.update", "permissions.manage", "payments.configure", "promotions.manage", "media.create", "media.settings.manage"}

func ValidPermission(key string) bool {
	for _, k := range PermissionKeys {
		if k == key {
			return true
		}
	}
	return false
}
func DefaultPermissionGrants(role UserRole) []PermissionGrant {
	grants := []PermissionGrant{{"offers.read", "own"}, {"rfqs.read", "own"}, {"orders.read", "own"}, {"orders.cancel", "own"}, {"rfqs.award", "own"}, {"rfqs.close", "own"}}
	if role == RoleAdmin {
		grants = []PermissionGrant{}
		for _, k := range PermissionKeys {
			grants = append(grants, PermissionGrant{k, "all"})
		}
		return grants
	}
	if role == RoleMerchant {
		for _, k := range []string{"workspace.access", "products.read", "products.create", "products.update", "products.delete", "offers.decide", "rfqs.quote", "payments.configure", "media.create"} {
			grants = append(grants, PermissionGrant{k, "own"})
		}
	}
	if role != RoleBuyer && role != RoleMerchant {
		return nil
	}
	return grants
}
func (u *User) Grants() []PermissionGrant {
	if u == nil || u.ID == "" {
		return nil
	}
	if u.PermissionsLoaded {
		return u.PermissionGrants
	}
	return DefaultPermissionGrants(u.Role)
}
func (u *User) Can(key, ownerID string) bool {
	if u == nil || u.ID == "" || !ValidPermission(key) {
		return false
	}
	for _, g := range u.Grants() {
		if g.Permission == key && (g.Scope == "all" || (g.Scope == "own" && ownerID != "" && ownerID == u.ID)) {
			return true
		}
	}
	return false
}
func (u *User) HasPermission(key string) bool {
	if u == nil {
		return false
	}
	return u.Can(key, u.ID)
}
func (u *User) CanAccessWorkspace() bool {
	for _, k := range []string{"workspace.access", "products.read", "products.create", "products.update", "products.delete", "offers.decide", "rfqs.quote", "users.read", "content.update", "brand.update", "permissions.manage", "payments.configure", "promotions.manage", "media.create", "media.settings.manage"} {
		if u.HasPermission(k) {
			return true
		}
	}
	return u.Can("orders.read", "") || u.Can("offers.read", "") || u.Can("rfqs.read", "")
}
