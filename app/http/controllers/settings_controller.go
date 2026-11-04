package controllers

import (
	"elcomware/app/facades"
	mw "elcomware/app/http/middleware"
	"elcomware/app/models"
	"elcomware/app/queries"
	"elcomware/app/services"
	"github.com/goravel/framework/contracts/http"
	inertia "github.com/goravel/inertia/facades"
	"net/url"
	"strings"
)

type SettingsController struct{}

func NewSettingsController() *SettingsController { return &SettingsController{} }
func (s *SettingsController) Index(ctx http.Context) http.Response {
	u := mw.CurrentUser(ctx)
	if u == nil {
		return ctx.Response().Redirect(302, "/login")
	}
	brand := u.Can("brand.update", "")
	access := u.Can("permissions.manage", "")
	paymentMerchant := ctx.Request().Query("merchant_id", u.ID)
	payments := u.Can("payments.configure", paymentMerchant)
	if paymentMerchant != u.ID {
		if !payments {
			return ctx.Response().String(403, "Payment configuration access denied")
		}
		var target models.User
		if facades.Orm().Query().Where("id = ?", paymentMerchant).First(&target) != nil || target.ID == "" {
			return ctx.Response().String(404, "Merchant account not found")
		}
	}
	if !brand && !access && !payments {
		return ctx.Response().String(403, "Settings access denied")
	}
	paymentPreference := services.DefaultPaymentPreference()
	if payments {
		paymentPreference = services.LoadPaymentPreference(paymentMerchant)
	}
	merchantChoices := []map[string]string{}
	if u.Can("payments.configure", "") {
		var accounts []models.User
		if err := facades.Orm().Query().OrderBy("name").Get(&accounts); err != nil {
			return ctx.Response().String(503, "Merchant choices unavailable")
		}
		for _, account := range accounts {
			queries.HydrateAccess(&account)
			if account.HasPermission("payments.configure") || account.ID == paymentMerchant {
				merchantChoices = append(merchantChoices, map[string]string{"id": account.ID, "name": account.Name, "company_name": account.CompanyName})
			}
		}
	}
	profiles := []map[string]any{}
	if access {
		var users []models.User
		if err := facades.Orm().Query().Get(&users); err != nil {
			return ctx.Response().String(500, "Settings unavailable")
		}
		for _, user := range users {
			queries.HydrateAccess(&user)
			version := 0
			var p models.AccessProfile
			if facades.Orm().Query().Where("user_id = ?", user.ID).First(&p) == nil {
				version = p.Version
			}
			profiles = append(profiles, map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "grants": user.Grants(), "version": version})
		}
	}
	return inertia.Inertia().Render(ctx, "Dashboard/Settings", map[string]any{"merchantChoices": merchantChoices, "currentUser": u, "canBrand": brand, "canPermissions": access, "canPayments": payments, "profiles": profiles, "permissionKeys": models.PermissionKeys, "brand": queries.LoadBrandSettings(), "paymentMerchantID": paymentMerchant, "canConfigureAllPayments": u.Can("payments.configure", ""), "paymentPreference": paymentPreference})
}
func settingsResult(ctx http.Context, err error, msg string) http.Response {
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", msg)
	}
	destination := "/dashboard/settings"
	if merchant := ctx.Request().Query("merchant_id", ""); merchant != "" {
		destination += "?merchant_id=" + url.QueryEscape(merchant)
	}
	if panel := ctx.Request().Query("panel", ""); panel == "brand" || panel == "access" || panel == "payments" {
		separator := "?"
		if strings.Contains(destination, "?") {
			separator = "&"
		}
		destination += separator + "panel=" + panel
	}
	return ctx.Response().Redirect(303, destination)
}
func (s *SettingsController) Brand(ctx http.Context) http.Response {
	var b queries.BrandSettings
	if err := ctx.Request().Bind(&b); err != nil {
		return ctx.Response().String(422, "Invalid settings")
	}
	return settingsResult(ctx, services.SaveBrandSettings(mw.CurrentUser(ctx), b), "Brand settings saved.")
}
func (s *SettingsController) Permissions(ctx http.Context) http.Response {
	var r struct {
		UserID  string                   `json:"user_id"`
		Grants  []models.PermissionGrant `json:"grants"`
		Version int                      `json:"version"`
	}
	if err := ctx.Request().Bind(&r); err != nil {
		return ctx.Response().String(422, "Invalid grants")
	}
	return settingsResult(ctx, services.SaveAccessProfile(mw.CurrentUser(ctx), r.UserID, r.Grants, r.Version), "Permissions saved and effective on the next request.")
}
func (s *SettingsController) Payments(ctx http.Context) http.Response {
	u := mw.CurrentUser(ctx)
	if u == nil {
		return ctx.Response().String(403, "Sign in required")
	}
	var p services.PaymentPreference
	if err := ctx.Request().Bind(&p); err != nil {
		return ctx.Response().String(422, "Invalid payment preferences")
	}
	merchant := ctx.Request().Query("merchant_id", u.ID)
	if !u.Can("payments.configure", merchant) {
		return ctx.Response().String(403, "Payment configuration access denied")
	}
	var target models.User
	if facades.Orm().Query().Where("id = ?", merchant).First(&target) != nil || target.ID == "" {
		return ctx.Response().String(404, "Merchant account not found")
	}
	err := services.SavePaymentPreference(u, merchant, p)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Payment preferences saved. Provider activation is still required.")
	}
	destination := "/dashboard/settings?merchant_id=" + url.QueryEscape(merchant)
	if ctx.Request().Query("panel", "") == "payments" {
		destination += "&panel=payments"
	}
	return ctx.Response().Redirect(303, destination)
}
