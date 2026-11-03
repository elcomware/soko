package actions

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/orm"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type ExecuteCheckoutAction struct {
	orderMapper *mappers.OrderMapper
}

func NewExecuteCheckoutAction() *ExecuteCheckoutAction {
	return &ExecuteCheckoutAction{
		orderMapper: mappers.NewOrderMapper(),
	}
}

func generateLicenseKey() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	hexStr := strings.ToUpper(hex.EncodeToString(bytes))
	return fmt.Sprintf("SOKO-%s-%s", hexStr[0:8], hexStr[8:16])
}

func (a *ExecuteCheckoutAction) Execute(input dto.CheckoutInput) (dto.OrderOutput, error) {
	if len(input.Items) == 0 {
		return dto.OrderOutput{}, fmt.Errorf("checkout cart cannot be empty")
	}

	orderID := fmt.Sprintf("ord_%d", time.Now().UnixNano())
	orderNumber := fmt.Sprintf("SOKO-ORD-%d", time.Now().Unix()%1000000)

	var totalAmount float64
	hasDigital := false
	hasPhysical := false

	orderItems := make([]models.OrderItem, 0, len(input.Items))

	for _, item := range input.Items {
		var product models.Product
		if err := facades.Orm().Query().Where("id = ?", item.ProductID).First(&product); err != nil {
			// Fallback placeholder if product was deleted
			product = models.Product{
				ID:    item.ProductID,
				Title: "Soko Marketplace Product",
				Type:  models.TypePhysical,
				Price: item.UnitPrice,
			}
		}

		unitPrice := item.UnitPrice
		if unitPrice <= 0 {
			unitPrice = product.Price
		}
		itemTotal := unitPrice * float64(item.Quantity)
		totalAmount += itemTotal

		isDigital := product.Type == models.TypeDigitalSoftware ||
			product.Type == models.TypeDigitalAsset ||
			product.Type == models.TypeDigitalLicense ||
			product.Format == models.FormatDigital

		if isDigital {
			hasDigital = true
		} else {
			hasPhysical = true
		}

		var licenseKey string
		var downloadUrl string
		if isDigital {
			licenseKey = generateLicenseKey()
			downloadUrl = product.DigitalDownloadUrl
			if downloadUrl == "" {
				downloadUrl = fmt.Sprintf("https://soko-storage.local/downloads/%s.zip", product.ID)
			}
		}

		orderItems = append(orderItems, models.OrderItem{
			ID:                 fmt.Sprintf("item_%d_%s", time.Now().UnixNano(), item.ProductID),
			OrderID:            orderID,
			ProductID:          product.ID,
			ProductTitle:       product.Title,
			ProductType:        product.Type,
			Quantity:           item.Quantity,
			UnitPrice:          unitPrice,
			TotalPrice:         itemTotal,
			DigitalLicenseKey:  licenseKey,
			DigitalDownloadUrl: downloadUrl,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		})
	}

	currency := input.Currency
	if currency == "" {
		currency = "USD"
	}
	paymentMethod := input.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "Trade Assurance Escrow"
	}

	order := models.Order{
		ID:            orderID,
		OrderNumber:   orderNumber,
		UserID:        input.UserID,
		CustomerName:  input.CustomerName,
		CustomerEmail: input.CustomerEmail,
		TotalAmount:   totalAmount,
		Currency:      currency,
		Status:        models.OrderStatusCompleted,
		PaymentMethod: paymentMethod,
		HasDigital:    hasDigital,
		HasPhysical:   hasPhysical,
		Items:         orderItems,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := facades.Orm().Transaction(func(tx orm.Query) error {
		if err := tx.Create(&order); err != nil {
			return err
		}
		for _, it := range orderItems {
			if err := tx.Create(&it); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return dto.OrderOutput{}, err
	}

	return a.orderMapper.ToDTO(order), nil
}
