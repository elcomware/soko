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
	return inertiafacades.Inertia().Render(ctx, "Dashboard/Promotions", map[string]any{"campaigns": result.Campaigns, "campaignPage": result.CampaignPage, "campaignCounts": result.CampaignCounts, "selectedCampaign": result.SelectedCampaign, "selectedCampaignStatus": result.SelectedCampaignStatus, "nextStatusChange": result.NextStatusChange, "products": result.Products, "canManageAll": actor.Can("promotions.manage", ""), "actorId": actor.ID, "serverTime": result.ServerTime})
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
