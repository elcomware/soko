package controllers

import (
	"elcomware/app/facades"
	appmiddleware "elcomware/app/http/middleware"
	"elcomware/app/models"
	"elcomware/app/services"
	"encoding/json"
	"errors"
	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"
	"net/url"
	"strconv"
	"time"
)

type PromotionController struct{}

func NewPromotionController() *PromotionController { return &PromotionController{} }
func (c *PromotionController) Index(ctx http.Context) http.Response {
	actor := appmiddleware.CurrentUser(ctx)
	if actor == nil || !actor.HasPermission("promotions.manage") {
		return ctx.Response().String(http.StatusForbidden, "Campaign access denied.")
	}
	input := map[string]string{}
	for _, key := range []string{"q", "status", "sort", "page", "per_page", "selected_campaign"} {
		input[key] = ctx.Request().Query(key, "")
	}
	filters, err := services.ParseWorkspacePromotionFilters(input)
	if err != nil {
		return ctx.Response().String(http.StatusBadRequest, "Invalid campaign filters: "+err.Error())
	}
	result, err := services.LoadWorkspacePromotions(actor, filters, time.Now().UTC())
	if err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "Campaign workspace could not be loaded. Please try again.")
	}
	return inertiafacades.Inertia().Render(ctx, "Dashboard/Promotions", map[string]any{"targetingApi": true, "campaigns": result.Campaigns, "campaignPage": result.CampaignPage, "campaignCounts": result.CampaignCounts, "selectedCampaign": result.SelectedCampaign, "selectedCampaignStatus": result.SelectedCampaignStatus, "nextStatusChange": result.NextStatusChange, "products": result.Products, "merchants": result.Merchants, "canManageAll": actor.Can("promotions.manage", ""), "actorId": actor.ID, "serverTime": result.ServerTime})
}
func (c *PromotionController) Save(ctx http.Context) http.Response {
	actor := appmiddleware.CurrentUser(ctx)
	if actor == nil || !actor.HasPermission("promotions.manage") {
		return ctx.Response().String(http.StatusForbidden, "Campaign access denied.")
	}
	p := models.Promotion{ID: ctx.Request().Route("id"), Name: ctx.Request().Input("name"), Kind: ctx.Request().Input("kind"), VendorID: ctx.Request().Input("vendor_id"), Enabled: ctx.Request().Input("enabled") == "true"}
	amount, e1 := strconv.ParseFloat(ctx.Request().Input("amount"), 64)
	start, e2 := time.Parse(time.RFC3339, ctx.Request().Input("starts_at"))
	end, e3 := time.Parse(time.RFC3339, ctx.Request().Input("ends_at"))
	e4 := json.Unmarshal([]byte(ctx.Request().Input("scope")), &p.Scope)
	var err error
	if !actor.Can("promotions.manage", p.VendorID) {
		return ctx.Response().String(http.StatusForbidden, "Campaign access denied.")
	}
	if p.VendorID != "" {
		var vendor models.User
		if lookupErr := facades.Orm().Query().Where("id = ?", p.VendorID).First(&vendor); lookupErr != nil || vendor.ID == "" {
			err = errors.New("select an existing merchant")
		} else if vendor.Role != models.RoleMerchant && vendor.Role != models.RoleAdmin {
			var owned models.Product
			if lookupErr := facades.Orm().Query().Where("vendor_id = ?", vendor.ID).First(&owned); lookupErr != nil || owned.ID == "" {
				err = errors.New("selected account does not own a catalog")
			}
		}
	}
	if err != nil {
		// Preserve validation errors; do not create a silently ineffective campaign.
	} else if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		err = errors.New("provide valid campaign targeting, discount and start/end times")
	} else {
		p.Amount = amount
		p.StartsAt = start.UTC()
		p.EndsAt = end.UTC()
		_, err = services.SavePromotion(p, actor, ctx.Request().Input("version"))
	}
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Campaign saved. It applies only during its scheduled period.")
	}
	destination := "/dashboard/promotions"
	returnInput := map[string]string{}
	hasReturn := false
	for _, key := range []string{"q", "status", "sort", "page", "per_page"} {
		returnInput[key] = ctx.Request().Input("return_" + key)
		hasReturn = hasReturn || returnInput[key] != ""
	}
	if hasReturn {
		if filters, err := services.ParseWorkspacePromotionFilters(returnInput); err == nil {
			values := url.Values{"q": {filters.Query}, "status": {filters.Status}, "sort": {filters.Sort}, "page": {strconv.Itoa(filters.Page)}, "per_page": {strconv.Itoa(filters.PerPage)}}
			destination += "?" + values.Encode()
		}
	}
	return ctx.Response().Redirect(http.StatusSeeOther, destination)
}

// Targets is a read-only JSON query; it does not save the draft scope.
func (c *PromotionController) Targets(ctx http.Context) http.Response {
	actor := appmiddleware.CurrentUser(ctx)
	vendorID := ctx.Request().Input("vendor_id")
	if actor == nil || !actor.Can("promotions.manage", vendorID) {
		return ctx.Response().Json(http.StatusForbidden, map[string]any{"error": "Campaign target access denied."})
	}
	raw := ctx.Request().Input("scope")
	if len(raw) > 1024*1024 {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{"error": "Target scope is too large."})
	}
	f := services.PromotionTargetFilters{VendorID: vendorID, Query: ctx.Request().Input("q"), Page: 1, PerPage: 25}
	if raw == "" {
		raw = "{}"
	}
	if err := json.Unmarshal([]byte(raw), &f.Scope); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{"error": "Invalid target scope."})
	}
	for key, target := range map[string]*int{"page": &f.Page, "per_page": &f.PerPage} {
		if raw := ctx.Request().Input(key); raw != "" {
			value, err := strconv.Atoi(raw)
			if err != nil {
				return ctx.Response().Json(http.StatusBadRequest, map[string]any{"error": "Invalid target pagination."})
			}
			*target = value
		}
	}
	if err := services.ValidatePromotionTargetFilters(f); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	result, err := services.LoadPromotionTargets(actor, f)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": "Campaign targets could not be loaded. Please try again."})
	}
	return ctx.Response().Json(http.StatusOK, result)
}
