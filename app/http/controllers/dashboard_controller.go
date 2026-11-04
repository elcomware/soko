package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	appmiddleware "elcomware/app/http/middleware"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
	"elcomware/app/services"
)

type DashboardController struct {
	dashboardService *services.DashboardService
	catalogService   *services.CatalogService
	orderService     *services.OrderService
	rfqService       *services.RfqService
	offerService     *services.OfferService
	currencyService  *services.CurrencyService
	authService      *services.AuthService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: services.NewDashboardService(),
		catalogService:   services.NewCatalogService(),
		orderService:     services.NewOrderService(),
		rfqService:       services.NewRfqService(),
		offerService:     services.NewOfferService(),
		currencyService:  services.NewCurrencyService(),
		authService:      services.NewAuthService(),
	}
}

func (d *DashboardController) getAuthenticatedUser(ctx http.Context) *models.User {
	return appmiddleware.CurrentUser(ctx)
}
func (d *DashboardController) Index(ctx http.Context) http.Response {
	u := d.getAuthenticatedUser(ctx)
	if u == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}
	if !u.CanManageCatalog() {
		return ctx.Response().Redirect(http.StatusFound, "/account")
	}
	data, err := services.LoadWorkspace(u, false)
	if err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "Your workspace could not be loaded. Please try again.")
	}
	return inertiafacades.Inertia().Render(ctx, "Dashboard/Index", map[string]any{"currentUser": u, "stats": data.Stats, "products": data.Products, "orders": data.Orders, "users": data.Users, "rfqs": data.Rfqs, "offers": data.Offers, "currencies": d.currencyService.GetCurrencies(), "regions": d.currencyService.GetRegions(), "capabilities": map[string]bool{"manageCatalog": u.CanManageCatalog(), "manageUsers": u.Role == models.RoleAdmin, "reviewOffers": u.CanManageCatalog(), "quoteRfqs": u.CanManageCatalog()}, "storefrontContent": queries.LoadStorefrontContent()})
}

func (d *DashboardController) CreateProduct(ctx http.Context) http.Response {
	var req requests.CreateProductRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil || errors != nil {
		errorMsg := "Invalid product input."
		if errors != nil {
			errorMsg = errors.One()
		}
		ctx.Request().Session().Flash("error", errorMsg)
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=products")
	}

	currentUser := d.getAuthenticatedUser(ctx)
	product, err := d.catalogService.CreateProduct(&req, currentUser)
	if err != nil {
		ctx.Request().Session().Flash("error", "Failed to create product: "+err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=products")
	}

	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise listing '%s' published to marketplace!", product.Title))
	return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=products")
}

func (d *DashboardController) DeleteProduct(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	if id != "" {
		if err := d.catalogService.DeleteProduct(id, d.getAuthenticatedUser(ctx)); err != nil {
			ctx.Request().Session().Flash("error", "Failed to delete product: "+err.Error())
		} else {
			ctx.Request().Session().Flash("success", "Product deleted from catalog.")
		}
	}
	return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=products")
}

func (d *DashboardController) SubmitRfqQuote(ctx http.Context) http.Response {
	var req requests.SubmitRfqQuoteRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil || errors != nil {
		errorMsg := "Quotation unit price is required."
		if errors != nil {
			errorMsg = errors.One()
		}
		ctx.Request().Session().Flash("error", errorMsg)
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=rfqs")
	}

	currentUser := d.getAuthenticatedUser(ctx)
	if _, err := d.rfqService.SubmitQuote(&req, currentUser); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=rfqs")
	}

	ctx.Request().Session().Flash("success", "Your competitive quotation has been submitted to the buyer!")
	return ctx.Response().Redirect(http.StatusFound, "/dashboard?section=rfqs")
}

func (d *DashboardController) UpdateProduct(ctx http.Context) http.Response {
	var req requests.CreateProductRequest
	if err := ctx.Request().Bind(&req); err != nil {
		ctx.Request().Session().Flash("error", "Provide valid product values.")
		return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard?section=products")
	}
	if _, err := d.catalogService.UpdateProduct(ctx.Request().Route("id"), &req, d.getAuthenticatedUser(ctx)); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Listing updated.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard?section=products")
}
func (d *DashboardController) DecideOffer(ctx http.Context) http.Response {
	if err := services.DecideOffer(ctx.Request().Route("id"), ctx.Request().Input("status"), d.getAuthenticatedUser(ctx)); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Offer decision saved. Payment has not been collected.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard?section=offers")
}
