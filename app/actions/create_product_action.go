package actions

import (
	"fmt"
	"strings"
	"time"

	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
)

type CreateProductAction struct{}

func NewCreateProductAction() *CreateProductAction {
	return &CreateProductAction{}
}

func (a *CreateProductAction) ExecuteFromRequest(req *requests.CreateProductRequest, user *models.User) (*models.Product, error) {
	price := req.Price
	if price <= 0 {
		price = 99.00
	}

	wholesalePrice := req.WholesalePrice
	if wholesalePrice <= 0 {
		wholesalePrice = price * 0.8
	}

	moq := req.MOQ
	if moq <= 0 {
		moq = 1
	}

	sku := strings.TrimSpace(req.Sku)
	if sku == "" {
		sku = fmt.Sprintf("SKU-SOK-%d", time.Now().Unix()%100000)
	}

	prodType := models.ProductType(req.Type)
	if prodType == "" {
		prodType = models.TypePhysical
	}

	format := models.ListingFormat(req.Format)
	if format == "" {
		format = models.FormatWholesale
	}

	condition := models.ItemCondition(req.Condition)
	if condition == "" {
		condition = models.ConditionNew
	}

	stock := req.Stock
	leadTime := 0
	incoterm := "Instant Cloud Delivery"
	digitalUrl := req.DigitalDownloadUrl
	digitalFormat := req.DigitalLicenseFormat
	digitalFileSize := req.DigitalFileSize

	if prodType == models.TypePhysical {
		if stock <= 0 {
			stock = 100
		}
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

	// Calculate Alibaba-style volume tiers automatically if wholesale
	var tiers []models.PriceTier
	if format == models.FormatWholesale {
		tiers = []models.PriceTier{
			{MinQty: moq, MaxQty: moq * 4, Price: price},
			{MinQty: moq*4 + 1, MaxQty: moq * 19, Price: price * 0.9},
			{MinQty: moq * 20, MaxQty: 9999, Price: wholesalePrice},
		}
	}

	vendorID := "usr_vendor_01"
	vendorName := "Soko Verified Supplier"
	vendorCountry := "Global"
	vendorBadge := "Gold Verified Supplier"

	if user != nil {
		vendorID = user.ID
		vendorName = user.CompanyName
		if vendorName == "" {
			vendorName = user.Name
		}
		if user.Country != "" {
			vendorCountry = user.Country
		}
		if user.Role == models.RoleAdmin {
			vendorBadge = "Platform Official"
		}
	}

	region := req.Region
	if region == "" {
		region = "global"
	}

	slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))
	product := &models.Product{
		ID:                   fmt.Sprintf("prod_%d", time.Now().UnixNano()),
		Title:                req.Title,
		Slug:                 fmt.Sprintf("%s-%d", slug, time.Now().Unix()%10000),
		Category:             req.Category,
		SubCategory:          req.SubCategory,
		Type:                 prodType,
		Format:               format,
		Condition:            condition,
		Price:                price,
		WholesalePrice:       wholesalePrice,
		PriceTiers:           tiers,
		MOQ:                  moq,
		Stock:                stock,
		Sku:                  sku,
		Description:          req.Description,
		ShortDesc:            req.Description,
		ImageUrl:             req.ImageUrl,
		VendorID:             vendorID,
		VendorName:           vendorName,
		VendorCountry:        vendorCountry,
		VendorRegion:         region,
		VendorVerified:       true,
		VendorBadge:          vendorBadge,
		VendorPositiveRate:   "99.5%",
		TradeAssuranceCover:  1000000.00,
		AcceptsOffers:        req.AcceptsOffers,
		Featured:             true,
		LeadTimeDays:         leadTime,
		Incoterm:             incoterm,
		DigitalDownloadUrl:   digitalUrl,
		DigitalLicenseFormat: digitalFormat,
		DigitalFileSize:      digitalFileSize,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	product.PrepareSave()
	if err := facades.Orm().Query().Create(product); err != nil {
		return nil, err
	}

	return product, nil
}
