package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

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
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user_id") != nil {
		if uid, ok := ctx.Request().Session().Get("user_id").(string); ok && uid != "" {
			if u, err := d.authService.GetUserByID(uid); err == nil && u != nil {
				return u
			}
		}
	}
	defaultUser, _ := d.authService.GetUserByID("usr_admin_01")
	return defaultUser
}

func (d *DashboardController) Index(ctx http.Context) http.Response {
	currentUser := d.getAuthenticatedUser(ctx)
	stats, _ := d.dashboardService.GetStats()
	products, _ := d.catalogService.GetProducts(queries.ProductFilter{})
	orders, _ := d.orderService.GetAllOrders()
	rfqs, _ := d.rfqService.GetOpenRfqs()
	offers, _ := d.offerService.GetOffers()
	currencies := d.currencyService.GetCurrencies()
	regions := d.currencyService.GetRegions()
	allUsers, _ := d.authService.GetAllUsers()

	return inertiafacades.Inertia().Render(ctx, "Dashboard/Index", map[string]any{
		"currentUser": currentUser,
		"stats":       stats,
		"products":    products,
		"orders":      orders,
		"users":       allUsers,
		"rfqs":        rfqs,
		"offers":      offers,
		"currencies":  currencies,
		"regions":     regions,
	})
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
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	currentUser := d.getAuthenticatedUser(ctx)
	product, err := d.catalogService.CreateProduct(&req, currentUser)
	if err != nil {
		ctx.Request().Session().Flash("error", "Failed to create product: "+err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise listing '%s' published to marketplace!", product.Title))
	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}

func (d *DashboardController) DeleteProduct(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	if id != "" {
		if err := d.catalogService.DeleteProduct(id); err != nil {
			ctx.Request().Session().Flash("error", "Failed to delete product: "+err.Error())
		} else {
			ctx.Request().Session().Flash("success", "Product deleted from catalog.")
		}
	}
	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
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
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	currentUser := d.getAuthenticatedUser(ctx)
	if _, err := d.rfqService.SubmitQuote(&req, currentUser); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	ctx.Request().Session().Flash("success", "Your competitive quotation has been submitted to the buyer!")
	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}
