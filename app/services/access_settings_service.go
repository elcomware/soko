package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/database/orm"
	"time"
)

func ValidateGrants(grants []models.PermissionGrant) error {
	if len(grants) > len(models.PermissionKeys) {
		return errors.New("too many grants")
	}
	seen := map[string]bool{}
	for _, g := range grants {
		if !models.ValidPermission(g.Permission) || (g.Scope != "own" && g.Scope != "all") || seen[g.Permission] {
			return errors.New("invalid or duplicate permission")
		}
		if g.Scope == "own" && (g.Permission == "permissions.manage" || g.Permission == "brand.update" || g.Permission == "content.update" || g.Permission == "users.read" || g.Permission == "media.settings.manage") {
			return errors.New("this permission requires platform scope")
		}
		seen[g.Permission] = true
	}
	return nil
}
func SaveAccessProfile(actor *models.User, target string, grants []models.PermissionGrant, version int) error {
	if !actor.Can("permissions.manage", "") {
		return errors.New("permission administration required")
	}
	if actor.ID == target {
		return errors.New("self permission changes are not allowed")
	}
	if err := ValidateGrants(grants); err != nil {
		return err
	}
	encoded, _ := json.Marshal(grants)
	return facades.Orm().Transaction(func(tx orm.Query) error {
		// Serialize permission changes before rechecking the actor, preventing mutually revoking administrators from racing.
		lock, err := tx.Exec("UPDATE access_mutex SET version = version + 1 WHERE id = ?", "permissions")
		if err != nil {
			return err
		}
		if lock.RowsAffected != 1 {
			return errors.New("permission storage is not ready")
		}
		var user models.User
		if err := tx.Where("id = ?", target).First(&user); err != nil {
			return err
		}
		if user.ID == "" {
			return errors.New("user not found")
		}
		var p models.AccessProfile
		if err := tx.Where("user_id = ?", target).First(&p); err != nil {
			return err
		}
		if p.Version != version {
			return errors.New("permissions changed; refresh before saving")
		}
		before := p.GrantsJSON
		if p.UserID == "" {
			b, _ := json.Marshal(models.DefaultPermissionGrants(user.Role))
			before = string(b)
		}
		// At least one other permission administrator remains: the freshly loaded actor.
		var actorProfile models.AccessProfile
		if err := tx.Where("user_id = ?", actor.ID).First(&actorProfile); err != nil {
			return err
		}
		fresh := *actor
		if actorProfile.UserID != "" {
			fresh.PermissionsLoaded = true
			fresh.PermissionGrants = nil
			if json.Unmarshal([]byte(actorProfile.GrantsJSON), &fresh.PermissionGrants) != nil {
				return errors.New("invalid actor profile")
			}
		}
		if !fresh.Can("permissions.manage", "") {
			return errors.New("permission revoked")
		}
		var previous []models.PermissionGrant
		if err := json.Unmarshal([]byte(before), &previous); err != nil {
			return errors.New("invalid target profile")
		}
		for _, g := range append(previous, grants...) {
			owner := fresh.ID
			if g.Scope == "all" {
				owner = ""
			}
			if !fresh.Can(g.Permission, owner) {
				return errors.New("you cannot delegate or change privileges beyond your own authority")
			}
		}

		if p.UserID == "" {
			p = models.AccessProfile{UserID: target, GrantsJSON: string(encoded), Version: 1}
			if err := tx.Create(&p); err != nil {
				return err
			}
		} else {
			r, err := tx.Model(&models.AccessProfile{}).Where("user_id = ? AND version = ?", target, version).Update(map[string]any{"grants_json": string(encoded), "version": version + 1})
			if err != nil {
				return err
			}
			if r.RowsAffected != 1 {
				return errors.New("permissions changed; refresh")
			}
		}
		return tx.Create(&models.AccessAudit{ID: uuid.NewString(), ActorID: actor.ID, TargetID: target, BeforeJSON: before, AfterJSON: string(encoded), CreatedAt: time.Now().Unix()})
	})
}
func SaveBrandSettings(actor *models.User, b queries.BrandSettings) error {
	if !actor.Can("brand.update", "") {
		return errors.New("brand update permission required")
	}
	if err := queries.ValidateBrandSettings(b); err != nil {
		return err
	}
	data, _ := json.Marshal(b)
	return facades.Orm().Transaction(func(tx orm.Query) error {
		var r models.BrandSettingsRecord
		if err := tx.Where("id = ?", "brand").First(&r); err != nil {
			return err
		}
		if r.Version != b.Version {
			return errors.New("brand settings changed; refresh before saving")
		}
		if r.ID == "" {
			return tx.Create(&models.BrandSettingsRecord{ID: "brand", SettingsJSON: string(data), Version: 1})
		}
		res, err := tx.Model(&models.BrandSettingsRecord{}).Where("id = ? AND version = ?", "brand", b.Version).Update(map[string]any{"settings_json": string(data), "version": b.Version + 1})
		if err == nil && res.RowsAffected != 1 {
			return errors.New("brand settings changed")
		}
		return err
	})
}
