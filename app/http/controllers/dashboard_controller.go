package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	"elcomware/app/models"
	"elcomware/app/services"
)

type DashboardController struct {
	store *services.Store
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		store: services.GetStore(),
	}
}

func (d *DashboardController) Index(ctx http.Context) http.Response {
	var currentUser *models.User
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
		if u, ok := ctx.Request().Session().Get("user").(*models.User); ok {
			currentUser = u
		}
	}

	if currentUser == nil {
		currentUser = d.store.GetUserByID("usr-admin-1")
	}

	stats := d.store.GetDashboardStats()
	products := d.store.GetProducts("all", "all", "all", "all", "all", "")
	orders := d.store.GetOrders()
	users := d.store.GetUsers()
	rfqs := d.store.GetRfqs()
	offers := d.store.GetOffers()
	currencies := d.store.GetCurrencies()
	regions := d.store.GetRegions()

	return inertiafacades.Inertia().Render(ctx, "Dashboard/Index", map[string]any{
		"currentUser": currentUser,
		"stats":       stats,
		"products":    products,
		"orders":      orders,
		"users":       users,
		"rfqs":        rfqs,
		"offers":      offers,
		"currencies":  currencies,
		"regions":     regions,
	})
}

func (d *DashboardController) CreateProduct(ctx http.Context) http.Response {
	title := strings.TrimSpace(ctx.Request().Input("title"))
	category := strings.TrimSpace(ctx.Request().Input("category"))
	subCategory := strings.TrimSpace(ctx.Request().Input("sub_category"))
	prodType := models.ProductType(strings.TrimSpace(ctx.Request().Input("type")))
	format := models.ListingFormat(strings.TrimSpace(ctx.Request().Input("format", "wholesale")))
	condition := models.ItemCondition(strings.TrimSpace(ctx.Request().Input("condition", "new")))
	region := strings.TrimSpace(ctx.Request().Input("region", "global"))

	priceStr := ctx.Request().Input("price", "99.00")
	price, _ := strconv.ParseFloat(priceStr, 64)
	if price <= 0 {
		price = 99.00
	}

	wholesaleStr := ctx.Request().Input("wholesale_price", "79.00")
	wholesalePrice, _ := strconv.ParseFloat(wholesaleStr, 64)
	if wholesalePrice <= 0 {
		wholesalePrice = price * 0.8
	}

	moq := ctx.Request().InputInt("moq", 1)
	stock := ctx.Request().InputInt("stock", 100)
	sku := strings.TrimSpace(ctx.Request().Input("sku"))
	description := strings.TrimSpace(ctx.Request().Input("description"))
	imageUrl := strings.TrimSpace(ctx.Request().Input("image_url"))
	digitalUrl := strings.TrimSpace(ctx.Request().Input("digital_download_url"))
	digitalFormat := strings.TrimSpace(ctx.Request().Input("digital_license_format"))
	digitalFileSize := strings.TrimSpace(ctx.Request().Input("digital_file_size"))
	acceptsOffers := ctx.Request().InputBool("accepts_offers", true)

	if title == "" {
		ctx.Request().Session().Flash("error", "Product title is required.")
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	if sku == "" {
		sku = fmt.Sprintf("SKU-SOK-%d", time.Now().Unix()%10000)
	}

	vendorName := "Soko Verified Supplier"
	vendorCountry := "Global"
	vendorID := "usr-merchant-1"
	vendorBadge := "Gold Verified Supplier"

	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
		if u, ok := ctx.Request().Session().Get("user").(*models.User); ok {
			vendorName = u.CompanyName
			if vendorName == "" {
				vendorName = u.Name
			}
			vendorCountry = u.Country
			vendorID = u.ID
			if u.Role == models.RoleAdmin {
				vendorBadge = "Platform Official"
			}
		}
	}

	if prodType == "" {
		prodType = models.TypePhysical
	}

	leadTime := 0
	incoterm := "Instant Cloud Delivery"
	if prodType == models.TypePhysical {
		leadTime = 7
		incoterm = "DDP Air Cargo / CIF"
	} else {
		stock = -1
		format = models.FormatDigital
		if digitalFormat == "" {
			digitalFormat = "SOKO-LIC-XXXX-XXXX-XXXX"
		}
		if digitalUrl == "" {
			digitalUrl = "https://cdn.soko.global/download/" + strings.ToLower(sku) + ".pkg"
		}
		if digitalFileSize == "" {
			digitalFileSize = "340 MB"
		}
	}

	// Build default Alibaba-style tiers if wholesale
	var tiers []models.PriceTier
	if format == models.FormatWholesale {
		tiers = []models.PriceTier{
			{MinQty: moq, MaxQty: moq * 4, Price: price},
			{MinQty: moq * 4 + 1, MaxQty: moq * 19, Price: price * 0.9},
			{MinQty: moq * 20, MaxQty: 9999, Price: wholesalePrice},
		}
	}

	newProduct := &models.Product{
		Title:               title,
		Slug:                strings.ToLower(strings.ReplaceAll(title, " ", "-")),
		Category:            category,
		SubCategory:         subCategory,
		Type:                prodType,
		Format:              format,
		Condition:           condition,
		Price:               price,
		WholesalePrice:      wholesalePrice,
		PriceTiers:          tiers,
		MOQ:                 moq,
		Stock:               stock,
		Sku:                 sku,
		Description:         description,
		ShortDesc:           description,
		ImageUrl:            imageUrl,
		VendorID:            vendorID,
		VendorName:          vendorName,
		VendorCountry:       vendorCountry,
		VendorRegion:        region,
		VendorVerified:      true,
		VendorBadge:         vendorBadge,
		VendorPositiveRate:  "99.5%",
		TradeAssuranceCover: 1000000.00,
		AcceptsOffers:       acceptsOffers,
		Rating:              5.0,
		ReviewsCount:        1,
		Featured:            true,
		LeadTimeDays:        leadTime,
		Incoterm:            incoterm,
		DigitalDownloadUrl:  digitalUrl,
		DigitalLicenseFormat: digitalFormat,
		DigitalFileSize:     digitalFileSize,
		CreatedAt:           time.Now(),
	}

	d.store.AddProduct(newProduct)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise listing '%s' published to %s marketplace!", title, strings.ToUpper(region)))
	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}

func (d *DashboardController) DeleteProduct(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	if id != "" {
		d.store.DeleteProduct(id)
		ctx.Request().Session().Flash("success", "Product deleted from catalog.")
	}
	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}

func (d *DashboardController) SubmitRfqQuote(ctx http.Context) http.Response {
	rfqID := ctx.Request().Input("rfq_id")
	unitPriceStr := ctx.Request().Input("unit_price")
	leadTime := ctx.Request().InputInt("lead_time_days", 7)
	notes := ctx.Request().Input("notes")

	unitPrice, _ := strconv.ParseFloat(unitPriceStr, 64)
	if unitPrice <= 0 || rfqID == "" {
		ctx.Request().Session().Flash("error", "Valid unit price is required.")
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	vendorName := "Apex Technologies"
	vendorCountry := "Germany"
	vendorID := "usr-merchant-1"

	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
		if u, ok := ctx.Request().Session().Get("user").(*models.User); ok {
			vendorName = u.CompanyName
			vendorCountry = u.Country
			vendorID = u.ID
		}
	}

	quote := models.RfqQuote{
		VendorID:     vendorID,
		VendorName:   vendorName,
		VendorCountry: vendorCountry,
		UnitPrice:    unitPrice,
		TotalAmount:  unitPrice * 100, // estimated
		LeadTimeDays: leadTime,
		Notes:        notes,
	}

	_, err := d.store.AddRfqQuote(rfqID, quote)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	ctx.Request().Session().Flash("success", "Your competitive quotation has been submitted to the buyer!")
	return ctx.Response().Redirect(http.StatusFound, "/dashboard")
}
