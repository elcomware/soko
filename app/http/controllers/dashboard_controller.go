package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	"elcomware/app/facades"
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
	if !u.CanAccessWorkspace() {
		return ctx.Response().Redirect(http.StatusFound, "/account")
	}

	var data *services.WorkspaceData
	var recordCapabilities map[string]map[string]map[string]bool
	var productPage *services.WorkspaceProductPageMeta
	productCategories := []string{}
	var selectedProduct *models.Product
	selectedProductStatus := "none"
	if ctx.Request().Query("section", "") == "products" {
		if !u.HasPermission("products.read") {
			return ctx.Response().String(http.StatusForbidden, "Product access is not permitted.")
		}
		input := map[string]string{}
		for _, key := range []string{"q", "category", "sort", "page", "per_page", "selected_product"} {
			input[key] = ctx.Request().Query(key, "")
		}
		filters, err := services.ParseWorkspaceProductFilters(input)
		if err != nil {
			return ctx.Response().String(http.StatusBadRequest, "Invalid product filters: "+err.Error())
		}
		result, err := services.LoadWorkspaceProducts(u, filters)
		if err != nil {
			return ctx.Response().String(http.StatusInternalServerError, "Your product workspace could not be loaded. Please try again.")
		}
		data = &services.WorkspaceData{Products: result.Products, Orders: []models.Order{}, Offers: []models.Offer{}, Rfqs: []models.Rfq{}, Users: []models.User{}, Stats: map[string]any{"total_products": result.ProductPage.Total}}
		recordCapabilities = result.RecordCapabilities
		productPage = &result.ProductPage
		productCategories = result.ProductCategories
		selectedProduct = result.SelectedProduct
		selectedProductStatus = result.SelectedProductStatus
	} else {
		var err error
		data, err = services.LoadWorkspace(u, false)
		if err != nil {
			return ctx.Response().String(http.StatusInternalServerError, "Your workspace could not be loaded. Please try again.")
		}
		recordCapabilities, err = dashboardRecordCapabilities(u, data)
		if err != nil {
			return ctx.Response().String(http.StatusInternalServerError, "Your workspace could not be loaded. Please try again.")
		}
	}
	return inertiafacades.Inertia().Render(ctx, "Dashboard/Index", map[string]any{"productPage": productPage, "productCategories": productCategories, "selectedProduct": selectedProduct, "selectedProductStatus": selectedProductStatus, "recordCapabilities": recordCapabilities, "currentUser": u, "stats": data.Stats, "products": data.Products, "orders": data.Orders, "users": data.Users, "rfqs": data.Rfqs, "offers": data.Offers, "currencies": d.currencyService.GetCurrencies(), "regions": d.currencyService.GetRegions(), "capabilities": map[string]bool{"manageCatalog": u.CanManageCatalog(), "createProducts": u.HasPermission("products.create"), "createMedia": u.HasPermission("media.create"), "editProducts": u.HasPermission("products.update"), "deleteProducts": u.HasPermission("products.delete"), "managePromotions": u.HasPermission("promotions.manage"), "manageUsers": u.Can("users.read", ""), "manageContent": u.HasPermission("content.update"), "reviewOffers": u.HasPermission("offers.decide"), "quoteRfqs": u.HasPermission("rfqs.quote"), "readProducts": u.HasPermission("products.read"), "readOrders": u.HasPermission("orders.read"), "readOffers": u.HasPermission("offers.read"), "readRfqs": u.HasPermission("rfqs.read") || u.HasPermission("rfqs.quote"), "manageSettings": u.HasPermission("brand.update") || u.HasPermission("permissions.manage") || u.HasPermission("payments.configure")}, "storefrontContent": queries.LoadStorefrontContent()})
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

// dashboardRecordCapabilities is presentation only. Every mutation still checks
// its effective grant and current record state in the corresponding service.
func dashboardRecordCapabilities(u *models.User, data *services.WorkspaceData) (map[string]map[string]map[string]bool, error) {
	result := map[string]map[string]map[string]bool{}
	for _, kind := range []string{"products", "orders", "offers", "rfqs", "quotes"} {
		result[kind] = map[string]map[string]bool{}
	}
	if len(data.Products) > 0 {
		ids := make([]string, 0, len(data.Products))
		for _, p := range data.Products {
			ids = append(ids, p.ID)
		}
		var referenced []models.Product
		// The deletion action preserves every listing with order or offer history.
		if err := facades.Orm().Query().Select("id").Where("id IN (?)", ids).Where("EXISTS (SELECT 1 FROM order_items WHERE product_id = products.id) OR EXISTS (SELECT 1 FROM offers WHERE product_id = products.id)").Get(&referenced); err != nil {
			return nil, err
		}
		history := map[string]bool{}
		for _, p := range referenced {
			history[p.ID] = true
		}
		for _, p := range data.Products {
			update := u.CanManageProduct(&p)
			result["products"][p.ID] = map[string]bool{"update": update, "delete": u.Can("products.delete", p.VendorID) && !history[p.ID], "media": update && u.Can("media.create", u.ID)}
		}
	}
	for _, order := range data.Orders {
		result["orders"][order.ID] = map[string]bool{"cancel": u.Can("orders.read", "") && order.Status == models.OrderStatusPending && u.Can("orders.cancel", order.UserID)}
	}
	// Resolve the offer's owner independently of visible catalog rows. An actor
	// can have offers.decide without products.read, and read-all/write-own grants.
	if len(data.Offers) > 0 {
		ids := make([]string, 0, len(data.Offers))
		for _, offer := range data.Offers {
			ids = append(ids, offer.ProductID)
		}
		var products []models.Product
		if err := facades.Orm().Query().Select("id", "vendor_id").Where("id IN (?)", ids).Get(&products); err != nil {
			return nil, err
		}
		owners := map[string]string{}
		for _, p := range products {
			owners[p.ID] = p.VendorID
		}
		for _, offer := range data.Offers {
			owner, exists := owners[offer.ProductID]
			result["offers"][offer.ID] = map[string]bool{"decide": exists && offer.Status == models.OfferStatusPending && u.Can("offers.decide", owner)}
		}
	}
	for _, rfq := range data.Rfqs {
		open := rfq.Status == models.RfqStatusOpen
		award := open && u.Can("rfqs.award", rfq.BuyerID)
		result["rfqs"][rfq.ID] = map[string]bool{"quote": open && rfq.BuyerID != u.ID && u.HasPermission("rfqs.quote"), "award": award, "close": open && u.Can("rfqs.close", rfq.BuyerID)}
		for _, quote := range rfq.Quotes {
			result["quotes"][quote.ID] = map[string]bool{"accept": award && quote.Status == "pending" && quote.RfqID == rfq.ID}
		}
	}
	return result, nil
}
