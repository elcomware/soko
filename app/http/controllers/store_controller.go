package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	"elcomware/app/models"
	"elcomware/app/services"
)

type StoreController struct {
	store *services.Store
}

func NewStoreController() *StoreController {
	return &StoreController{
		store: services.GetStore(),
	}
}

func (s *StoreController) Index(ctx http.Context) http.Response {
	filterType := ctx.Request().Query("type", "all")
	format := ctx.Request().Query("format", "all")
	region := ctx.Request().Query("region", "global")
	category := ctx.Request().Query("category", "all")
	condition := ctx.Request().Query("condition", "all")
	query := ctx.Request().Query("q", "")

	products := s.store.GetProducts(filterType, format, region, category, condition, query)
	stats := s.store.GetDashboardStats()
	currencies := s.store.GetCurrencies()
	regions := s.store.GetRegions()
	rfqs := s.store.GetRfqs()

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
		"products":        products,
		"categories":      categories,
		"currencies":      currencies,
		"regions":         regions,
		"rfqs":            rfqs,
		"activeType":      filterType,
		"activeFormat":    format,
		"activeRegion":    region,
		"activeCategory":  category,
		"activeCondition": condition,
		"searchQuery":     query,
		"stats":           stats,
	})
}

func (s *StoreController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	product := s.store.GetProductByID(id)
	if product == nil {
		ctx.Request().Session().Flash("error", "Product not found.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	related := s.store.GetProducts(string(product.Type), "", "", product.Category, "", "")
	filteredRelated := make([]*models.Product, 0)
	for _, p := range related {
		if p.ID != product.ID && len(filteredRelated) < 3 {
			filteredRelated = append(filteredRelated, p)
		}
	}

	currencies := s.store.GetCurrencies()
	regions := s.store.GetRegions()

	return inertiafacades.Inertia().Render(ctx, "Store/ProductDetail", map[string]any{
		"product":    product,
		"related":    filteredRelated,
		"currencies": currencies,
		"regions":    regions,
	})
}

func (s *StoreController) Checkout(ctx http.Context) http.Response {
	productID := ctx.Request().Input("product_id")
	quantity := ctx.Request().InputInt("quantity", 1)
	customerName := ctx.Request().Input("customer_name")
	customerEmail := ctx.Request().Input("customer_email")
	currencyCode := ctx.Request().Input("currency", "USD")
	paymentRail := ctx.Request().Input("payment_method", "Trade Escrow Guarantee")

	product := s.store.GetProductByID(productID)
	if product == nil {
		ctx.Request().Session().Flash("error", "Invalid product selected for checkout.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	if quantity < product.MOQ {
		quantity = product.MOQ
	}

	// Calculate price according to Alibaba volume tier if available
	unitPrice := product.Price
	if len(product.PriceTiers) > 0 {
		for _, tier := range product.PriceTiers {
			if quantity >= tier.MinQty && (tier.MaxQty == 0 || quantity <= tier.MaxQty) {
				unitPrice = tier.Price
				break
			}
		}
	} else if quantity >= 10 && product.WholesalePrice > 0 {
		unitPrice = product.WholesalePrice
	}

	total := float64(quantity) * unitPrice

	userID := "guest"
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
		if u, ok := ctx.Request().Session().Get("user").(*models.User); ok {
			userID = u.ID
			if customerName == "" {
				customerName = u.Name
			}
			if customerEmail == "" {
				customerEmail = u.Email
			}
		}
	}

	if customerName == "" {
		customerName = "Enterprise Procurement Partner"
	}
	if customerEmail == "" {
		customerEmail = "procurement@globalbuyer.com"
	}

	order := &models.Order{
		UserID:        userID,
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		TotalAmount:   total,
		Currency:      currencyCode,
		PaymentMethod: paymentRail,
		Items: []models.OrderItem{
			{
				ProductID:    product.ID,
				ProductTitle: product.Title,
				ProductType:  product.Type,
				Quantity:     quantity,
				UnitPrice:    unitPrice,
				TotalPrice:   total,
			},
		},
	}

	createdOrder := s.store.CreateOrder(order)
	if createdOrder.HasDigital {
		ctx.Request().Session().Flash("success", "Escrow Verified! Automated license key issued & instant download token ready.")
	} else {
		ctx.Request().Session().Flash("success", fmt.Sprintf("Wholesale order %s placed under Soko Trade Escrow! Supplier dispatched notification.", createdOrder.OrderNumber))
	}

	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}

// Make an Offer (eBay Style)
func (s *StoreController) MakeOffer(ctx http.Context) http.Response {
	productID := ctx.Request().Input("product_id")
	offeredPriceStr := ctx.Request().Input("offered_price")
	quantity := ctx.Request().InputInt("quantity", 1)
	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	company := ctx.Request().Input("company")
	message := ctx.Request().Input("message")

	offeredPrice, _ := strconv.ParseFloat(offeredPriceStr, 64)
	if offeredPrice <= 0 || productID == "" {
		ctx.Request().Session().Flash("error", "Valid offer price is required.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	product := s.store.GetProductByID(productID)
	if product == nil {
		ctx.Request().Session().Flash("error", "Product not found.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	offer := &models.Offer{
		ProductID:    product.ID,
		ProductTitle: product.Title,
		BuyerName:    name,
		BuyerEmail:   email,
		BuyerCompany: company,
		OfferedPrice: offeredPrice,
		Quantity:     quantity,
		Message:      message,
	}

	s.store.SubmitOffer(offer)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Your formal offer of $%.2f/unit for %s has been submitted to the supplier for review!", offeredPrice, product.Title))
	return ctx.Response().Redirect(http.StatusFound, "/")
}

// Request for Quotation (Alibaba Style)
func (s *StoreController) CreateRfq(ctx http.Context) http.Response {
	title := strings.TrimSpace(ctx.Request().Input("title"))
	category := strings.TrimSpace(ctx.Request().Input("category"))
	quantity := ctx.Request().InputInt("quantity", 10)
	targetPriceStr := ctx.Request().Input("target_price", "0")
	targetPrice, _ := strconv.ParseFloat(targetPriceStr, 64)
	destination := ctx.Request().Input("destination_country")
	specs := ctx.Request().Input("specifications")
	buyerName := ctx.Request().Input("buyer_name")
	buyerCompany := ctx.Request().Input("buyer_company")
	buyerCountry := ctx.Request().Input("buyer_country")

	if title == "" {
		ctx.Request().Session().Flash("error", "RFQ title is required.")
		return ctx.Response().Redirect(http.StatusFound, "/#rfq")
	}

	rfq := &models.Rfq{
		Title:              title,
		Category:           category,
		Quantity:           quantity,
		TargetPrice:        targetPrice,
		DestinationCountry: destination,
		Specifications:     specs,
		BuyerName:          buyerName,
		BuyerCompany:       buyerCompany,
		BuyerCountry:       buyerCountry,
	}

	created := s.store.CreateRfq(rfq)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise RFQ %s published to global suppliers! You will receive quotations within 24 hours.", created.RfqNumber))
	return ctx.Response().Redirect(http.StatusFound, "/#rfq")
}
