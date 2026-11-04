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

type StoreController struct {
	catalogService   *services.CatalogService
	orderService     *services.OrderService
	rfqService       *services.RfqService
	offerService     *services.OfferService
	currencyService  *services.CurrencyService
	dashboardService *services.DashboardService
	authService      *services.AuthService
}

func NewStoreController() *StoreController {
	return &StoreController{
		catalogService:   services.NewCatalogService(),
		orderService:     services.NewOrderService(),
		rfqService:       services.NewRfqService(),
		offerService:     services.NewOfferService(),
		currencyService:  services.NewCurrencyService(),
		dashboardService: services.NewDashboardService(),
		authService:      services.NewAuthService(),
	}
}

func (s *StoreController) getAuthenticatedUser(ctx http.Context) *models.User {
	return appmiddleware.CurrentUser(ctx)
}

func (s *StoreController) Index(ctx http.Context) http.Response {
	filter := queries.ProductFilter{
		Type:      ctx.Request().Query("type", "all"),
		Format:    ctx.Request().Query("format", "all"),
		Region:    models.NormalizeRegion(ctx.Request().Query("region", "global")),
		Category:  ctx.Request().Query("category", "all"),
		Condition: ctx.Request().Query("condition", "all"),
		Query:     ctx.Request().Query("q", ""),
	}

	products, err := s.catalogService.GetProducts(filter)
	if err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "The catalog could not be loaded. Please try again.")
	}
	products = models.PublicProducts(products)
	stats := map[string]any{"total_products": len(products)}
	currencies := s.currencyService.GetCurrencies()
	regions := s.currencyService.GetRegions()
	rfqs, _ := s.rfqService.GetOpenRfqs()

	categories := []map[string]any{
		{"id": "all", "name": "All Categories", "icon": "Layers", "count": len(products)},
		{"id": "Industrial & Machinery", "name": "Industrial & Heavy CNC", "type": "physical", "icon": "Tool"},
		{"id": "Electronics & Telecom", "name": "Datacenter & Networking", "type": "physical", "icon": "Wifi"},
		{"id": "Digital Software", "name": "Cloud Software & SaaS", "type": "digital", "icon": "Cloud"},
		{"id": "Digital Assets", "name": "AI Models & CAD Blueprints", "type": "digital", "icon": "Cpu"},
		{"id": "Agriculture & Processing", "name": "Agri-Tech & Sorting Mills", "type": "physical", "icon": "Leaf"},
		{"id": "Energy & Utilities", "name": "Commercial Energy & BESS", "type": "physical", "icon": "Zap"},
	}

	return inertiafacades.Inertia().Render(ctx, "Store/Index", map[string]any{
		"products":          products,
		"storefrontContent": queries.LoadStorefrontContent(),
		"categories":        categories,
		"currencies":        currencies,
		"regions":           regions,
		"rfqs":              rfqs,
		"activeType":        filter.Type,
		"activeFormat":      filter.Format,
		"activeRegion":      filter.Region,
		"activeCategory":    filter.Category,
		"activeCondition":   filter.Condition,
		"searchQuery":       filter.Query,
		"stats":             stats,
	})
}

func (s *StoreController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	product, err := s.catalogService.GetProductByID(id)
	if err != nil || product == nil || product.ID == "" {
		ctx.Request().Session().Flash("error", "Product not found.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	relatedList, _ := s.catalogService.GetProducts(queries.ProductFilter{
		Type:     string(product.Type),
		Category: product.Category,
	})

	filteredRelated := make([]models.Product, 0)
	for _, p := range relatedList {
		if p.ID != product.ID && len(filteredRelated) < 3 {
			filteredRelated = append(filteredRelated, p)
		}
	}

	currencies := s.currencyService.GetCurrencies()
	regions := s.currencyService.GetRegions()

	return inertiafacades.Inertia().Render(ctx, "Store/ProductDetail", map[string]any{
		"product":    product.Public(),
		"related":    models.PublicProducts(filteredRelated),
		"currencies": currencies,
		"regions":    regions,
	})
}

func (s *StoreController) Checkout(ctx http.Context) http.Response {
	return NewCommerceController().Checkout(ctx)
}

// Make an Offer (eBay Style)
func (s *StoreController) MakeOffer(ctx http.Context) http.Response {
	var req requests.SubmitOfferRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil || errors != nil {
		errorMsg := "Valid offer details are required."
		if errors != nil {
			errorMsg = errors.One()
		}
		ctx.Request().Session().Flash("error", errorMsg)
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	currentUser := s.getAuthenticatedUser(ctx)
	offer, err := s.offerService.SubmitOffer(&req, currentUser)
	if err != nil {
		ctx.Request().Session().Flash("error", "Failed to submit offer: "+err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	ctx.Request().Session().Flash("success", fmt.Sprintf("Your formal offer of $%.2f/unit for %s has been submitted to the supplier for review!", offer.OfferedPrice, offer.ProductTitle))
	return ctx.Response().Redirect(http.StatusFound, "/")
}

// Request for Quotation (Alibaba Style)
func (s *StoreController) CreateRfq(ctx http.Context) http.Response {
	var req requests.CreateRfqRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil || errors != nil {
		errorMsg := "RFQ brief details are required."
		if errors != nil {
			errorMsg = errors.One()
		}
		ctx.Request().Session().Flash("error", errorMsg)
		return ctx.Response().Redirect(http.StatusFound, "/#rfq")
	}

	currentUser := s.getAuthenticatedUser(ctx)
	rfq, err := s.rfqService.CreateRfq(&req, currentUser)
	if err != nil {
		ctx.Request().Session().Flash("error", "Failed to create RFQ: "+err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/#rfq")
	}

	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise RFQ %s published to global suppliers! Suppliers can now submit quotations.", rfq.RfqNumber))
	return ctx.Response().Redirect(http.StatusFound, "/#rfq")
}
