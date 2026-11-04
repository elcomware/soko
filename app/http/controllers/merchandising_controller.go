package controllers

import (
	"elcomware/app/facades"
	mw "elcomware/app/http/middleware"
	"elcomware/app/models"
	"elcomware/app/queries"
	"elcomware/app/services"
	"github.com/goravel/framework/contracts/http"
	inertia "github.com/goravel/inertia/facades"
)

type MerchandisingController struct{}

func (*MerchandisingController) Index(ctx http.Context) http.Response {
	actor := mw.CurrentUser(ctx)
	if !actor.Can("content.update", "") {
		return ctx.Response().String(403, "Storefront content permission required")
	}
	config, err := queries.LoadMerchandising()
	if err != nil {
		return ctx.Response().String(503, "Storefront settings storage is unavailable")
	}
	var products []models.Product
	if err := facades.Orm().Query().OrderBy("title").Get(&products); err != nil {
		return ctx.Response().String(503, "Product choices are unavailable")
	}
	choices := []map[string]any{}
	for _, p := range products {
		choices = append(choices, map[string]any{"id": p.ID, "title": p.Title, "category": p.Category, "digital": p.IsDigital(), "stock": p.Stock})
	}
	return inertia.Inertia().Render(ctx, "Dashboard/Merchandising", map[string]any{"settings": config, "products": choices})
}
func (*MerchandisingController) Save(ctx http.Context) http.Response {
	actor := mw.CurrentUser(ctx)
	if !actor.Can("content.update", "") {
		return ctx.Response().String(403, "Storefront content permission required")
	}
	var settings models.MerchandisingSettings
	if ctx.Request().Bind(&settings) != nil {
		return ctx.Response().String(422, "Invalid storefront settings")
	}
	if err := services.SaveMerchandising(actor, settings); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Storefront selection and sections published.")
	}
	return ctx.Response().Redirect(303, "/dashboard/merchandising")
}
