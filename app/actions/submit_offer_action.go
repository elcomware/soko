package actions

import (
	"errors"
	"fmt"
	"github.com/goravel/framework/contracts/database/orm"
	"math"
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
	var result *models.Offer
	err := facades.Orm().Transaction(func(tx orm.Query) error { var err error; result, err = a.executeFromQuery(req, user, tx); return err })
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (a *SubmitOfferAction) executeFromQuery(req *requests.SubmitOfferRequest, user *models.User, tx orm.Query) (*models.Offer, error) {
	var product models.Product
	if err := tx.Where("id = ?", req.ProductID).First(&product); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if user == nil || user.ID == "" {
		return nil, errors.New("sign in to make an offer")
	}
	if product.ID == "" || !product.AcceptsOffers || product.VendorID == user.ID {
		return nil, errors.New("this listing does not accept your offer")
	}
	if req.OfferedPrice*float64(req.Quantity) > 1e12 || req.OfferedPrice <= 0 || req.OfferedPrice > 1e12 || math.IsNaN(req.OfferedPrice) || math.IsInf(req.OfferedPrice, 0) {
		return nil, errors.New("offer price must be a positive finite amount")
	}
	if _, err := product.PriceFor(req.Quantity); err != nil {
		return nil, err
	}
	if len(req.Message) > 5000 {
		return nil, errors.New("offer message is too long")
	}
	buyerID := user.ID
	buyerName := user.Name
	buyerEmail := user.Email
	buyerCompany := user.CompanyName

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

	if err := tx.Create(offer); err != nil {
		return nil, err
	}

	return offer, nil
}
