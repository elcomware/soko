package actions

import (
	"fmt"
	"time"

	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
)

type SubmitOfferAction struct{}

func NewSubmitOfferAction() *SubmitOfferAction {
	return &SubmitOfferAction{}
}

func (a *SubmitOfferAction) ExecuteFromRequest(req *requests.SubmitOfferRequest, user *models.User) (*models.Offer, error) {
	var product models.Product
	if err := facades.Orm().Query().Where("id = ?", req.ProductID).First(&product); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	buyerID := "guest"
	buyerName := req.Name
	buyerEmail := req.Email
	buyerCompany := req.Company

	if user != nil {
		buyerID = user.ID
		if buyerName == "" {
			buyerName = user.Name
		}
		if buyerEmail == "" {
			buyerEmail = user.Email
		}
		if buyerCompany == "" {
			buyerCompany = user.CompanyName
		}
	}

	qty := req.Quantity
	if qty <= 0 {
		qty = 1
	}

	offer := &models.Offer{
		ID:           fmt.Sprintf("off_%d", time.Now().UnixNano()),
		ProductID:    product.ID,
		ProductTitle: product.Title,
		BuyerID:      buyerID,
		BuyerName:    buyerName,
		BuyerEmail:   buyerEmail,
		BuyerCompany: buyerCompany,
		OfferedPrice: req.OfferedPrice,
		Quantity:     qty,
		TotalOffered: req.OfferedPrice * float64(qty),
		Message:      req.Message,
		Status:       models.OfferStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := facades.Orm().Query().Create(offer); err != nil {
		return nil, err
	}

	return offer, nil
}
