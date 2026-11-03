package actions

import (
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/database/orm"

	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
)

type ExecuteCheckoutAction struct {
	licenseAction *IssueDigitalLicenseAction
}

func NewExecuteCheckoutAction() *ExecuteCheckoutAction {
	return &ExecuteCheckoutAction{
		licenseAction: NewIssueDigitalLicenseAction(),
	}
}

func (a *ExecuteCheckoutAction) ExecuteFromRequest(req *requests.CheckoutRequest, user *models.User) (*models.Order, error) {
	var product models.Product
	if err := facades.Orm().Query().Where("id = ?", req.ProductID).First(&product); err != nil {
		return nil, fmt.Errorf("invalid product selected for checkout: %w", err)
	}
	product.Hydrate()

	quantity := req.Quantity
	if quantity < product.MOQ {
		quantity = product.MOQ
	}

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

	totalAmount := float64(quantity) * unitPrice

	userID := "guest"
	customerName := req.CustomerName
	customerEmail := req.CustomerEmail

	if user != nil {
		userID = user.ID
		if customerName == "" {
			customerName = user.Name
		}
		if customerEmail == "" {
			customerEmail = user.Email
		}
	}

	if customerName == "" {
		customerName = "Enterprise Procurement Partner"
	}
	if customerEmail == "" {
		customerEmail = "procurement@globalbuyer.com"
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	paymentRail := req.PaymentMethod
	if paymentRail == "" {
		paymentRail = "Trade Assurance Escrow"
	}

	orderID := fmt.Sprintf("ord_%d", time.Now().UnixNano())
	orderNumber := fmt.Sprintf("SOKO-ORD-%d", time.Now().Unix()%1000000)

	isDigital := product.Type == models.TypeDigitalSoftware ||
		product.Type == models.TypeDigitalAsset ||
		product.Type == models.TypeDigitalLicense ||
		product.Format == models.FormatDigital

	var licenseKey, downloadUrl string
	if isDigital {
		licenseKey, downloadUrl = a.licenseAction.Execute(&product)
	}

	item := models.OrderItem{
		ID:                 fmt.Sprintf("item_%d", time.Now().UnixNano()),
		OrderID:            orderID,
		ProductID:          product.ID,
		ProductTitle:       product.Title,
		ProductType:        product.Type,
		Quantity:           quantity,
		UnitPrice:          unitPrice,
		TotalPrice:         totalAmount,
		DigitalLicenseKey:  licenseKey,
		DigitalDownloadUrl: downloadUrl,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	order := &models.Order{
		ID:            orderID,
		OrderNumber:   orderNumber,
		UserID:        userID,
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		TotalAmount:   totalAmount,
		Currency:      currency,
		Status:        models.OrderStatusCompleted,
		PaymentMethod: paymentRail,
		HasDigital:    isDigital,
		HasPhysical:   !isDigital,
		Items:         []models.OrderItem{item},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := facades.Orm().Transaction(func(tx orm.Query) error {
		if err := tx.Create(order); err != nil {
			return err
		}
		return tx.Create(&item)
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}
