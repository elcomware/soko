package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
)

type ContentController struct{}

func NewContentController() *ContentController { return &ContentController{} }

func (c *ContentController) Save(ctx http.Context) http.Response {
	if !ctx.Request().HasSession() {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}
	uid := ctx.Request().Session().Get("user_id")
	var actor models.User
	if uid == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}
	if err := facades.Orm().Query().Where("id = ?", fmt.Sprint(uid)).First(&actor); err != nil || actor.ID == "" || actor.Role != models.RoleAdmin {
		return ctx.Response().String(http.StatusForbidden, "Only administrators can update storefront content.")
	}
	values := make(map[string]string, 6)
	for key := range queries.DefaultStorefrontContent() {
		values[key] = ctx.Request().Input(key)
	}
	valid, err := queries.ValidateStorefrontContent(values)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=content")
	}
	if !facades.Schema().HasTable("storefront_contents") {
		ctx.Request().Session().Flash("error", "Storefront content storage is not ready. Apply the pending migration before saving.")
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=content")
	}
	encoded, err := json.Marshal(valid)
	if err == nil {
		err = facades.Orm().Transaction(func(tx orm.Query) error {
			var content models.StorefrontContent
			return tx.UpdateOrCreate(&content, models.StorefrontContent{ID: "store"}, map[string]any{"content_json": string(encoded)})
		})
	}
	if err != nil {
		ctx.Request().Session().Flash("error", "Storefront content could not be saved. Your changes have not been published.")
	} else {
		ctx.Request().Session().Flash("success", "Storefront content saved and published.")
	}
	return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=content")
}
