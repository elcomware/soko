package controllers

import (
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
	category := ctx.Request().Query("category", "all")
	query := ctx.Request().Query("q", "")

	products := s.store.GetProducts(filterType, category, query)
	stats := s.store.GetDashboardStats()

	categories := []map[string]any{
		{"id": "all", "name": "All Categories", "count": len(products), "icon": "Layers"},
		{"id": "Digital Software", "name": "Cloud Software & SaaS", "type": "digital", "icon": "Cloud"},
		{"id": "Digital Assets", "name": "AI Models & CAD Blueprints", "type": "digital", "icon": "Cpu"},
		{"id": "Industrial & Machinery", "name": "Industrial & Heavy CNC", "type": "physical", "icon": "Tool"},
		{"id": "Electronics & Telecom", "name": "Enterprise Networking & Telecom", "type": "physical", "icon": "Wifi"},
		{"id": "Energy & Utilities", "name": "Commercial Energy Storage", "type": "physical", "icon": "Zap"},
	}

	return inertiafacades.Inertia().Render(ctx, "Store/Index", map[string]any{
		"products":       products,
		"categories":     categories,
		"activeType":     filterType,
		"activeCategory": category,
		"searchQuery":    query,
		"stats":          stats,
	})
}

func (s *StoreController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	product := s.store.GetProductByID(id)
	if product == nil {
		ctx.Request().Session().Flash("error", "Product not found.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	related := s.store.GetProducts(string(product.Type), "", "")
	// Filter out the current product
	filteredRelated := make([]*models.Product, 0)
	for _, p := range related {
		if p.ID != product.ID && len(filteredRelated) < 3 {
			filteredRelated = append(filteredRelated, p)
		}
	}

	return inertiafacades.Inertia().Render(ctx, "Store/ProductDetail", map[string]any{
		"product": product,
		"related": filteredRelated,
	})
}

func (s *StoreController) Checkout(ctx http.Context) http.Response {
	productID := ctx.Request().Input("product_id")
	quantity := ctx.Request().InputInt("quantity", 1)
	customerName := ctx.Request().Input("customer_name")
	customerEmail := ctx.Request().Input("customer_email")

	product := s.store.GetProductByID(productID)
	if product == nil {
		ctx.Request().Session().Flash("error", "Invalid product selected for checkout.")
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	if quantity < product.MOQ {
		quantity = product.MOQ
	}

	unitPrice := product.Price
	if quantity >= 10 && product.WholesalePrice > 0 {
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
		Currency:      "USD",
		PaymentMethod: "Corporate Wire / Escrow",
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
		ctx.Request().Session().Flash("success", "Order fulfilled! Instant digital license and download access generated.")
	} else {
		ctx.Request().Session().Flash("success", "Wholesale order placed successfully! Supplier notified for freight dispatch.")
	}

	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}
