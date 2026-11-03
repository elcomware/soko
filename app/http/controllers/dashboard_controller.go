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

	// Fallback to demo admin if not directly typed in session
	if currentUser == nil {
		currentUser = d.store.GetUserByID("usr-admin-1")
	}

	stats := d.store.GetDashboardStats()
	products := d.store.GetProducts("all", "all", "")
	orders := d.store.GetOrders()
	users := d.store.GetUsers()

	return inertiafacades.Inertia().Render(ctx, "Dashboard/Index", map[string]any{
		"currentUser": currentUser,
		"stats":       stats,
		"products":    products,
		"orders":      orders,
		"users":       users,
	})
}

func (d *DashboardController) CreateProduct(ctx http.Context) http.Response {
	title := strings.TrimSpace(ctx.Request().Input("title"))
	category := strings.TrimSpace(ctx.Request().Input("category"))
	subCategory := strings.TrimSpace(ctx.Request().Input("sub_category"))
	prodType := models.ProductType(strings.TrimSpace(ctx.Request().Input("type")))

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

	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
		if u, ok := ctx.Request().Session().Get("user").(*models.User); ok {
			vendorName = u.CompanyName
			if vendorName == "" {
				vendorName = u.Name
			}
			vendorCountry = u.Country
			vendorID = u.ID
		}
	}

	if prodType == "" {
		prodType = models.TypeDigitalSoftware
	}

	leadTime := 0
	if prodType == models.TypePhysical {
		leadTime = 7
	} else {
		stock = -1 // unlimited for digital
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

	newProduct := &models.Product{
		Title:                title,
		Slug:                 strings.ToLower(strings.ReplaceAll(title, " ", "-")),
		Category:             category,
		SubCategory:          subCategory,
		Type:                 prodType,
		Price:                price,
		WholesalePrice:       wholesalePrice,
		MOQ:                  moq,
		Stock:                stock,
		Sku:                  sku,
		Description:          description,
		ShortDesc:            description,
		ImageUrl:             imageUrl,
		VendorID:             vendorID,
		VendorName:           vendorName,
		VendorCountry:        vendorCountry,
		VendorVerified:       true,
		Rating:               5.0,
		ReviewsCount:         1,
		Featured:             true,
		LeadTimeDays:         leadTime,
		DigitalDownloadUrl:   digitalUrl,
		DigitalLicenseFormat: digitalFormat,
		DigitalFileSize:      digitalFileSize,
		CreatedAt:            time.Now(),
	}

	d.store.AddProduct(newProduct)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Product '%s' added to enterprise catalog successfully!", title))
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
