package actions

import (
	"errors"
	"fmt"
	"math"
	"time"

	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
)

type CreateRfqAction struct{}

func NewCreateRfqAction() *CreateRfqAction {
	return &CreateRfqAction{}
}

func (a *CreateRfqAction) ExecuteFromRequest(req *requests.CreateRfqRequest, user *models.User) (*models.Rfq, error) {
	if user == nil || user.ID == "" {
		return nil, errors.New("sign in to publish a sourcing request")
	}
	if req.Quantity < 1 || req.Quantity > 1000000 || req.TargetPrice < 0 || req.TargetPrice > 1e12 || math.IsNaN(req.TargetPrice) || math.IsInf(req.TargetPrice, 0) || len(req.Title) < 5 || len(req.Title) > 255 || len(req.Specifications) > 20000 {
		return nil, errors.New("provide a valid sourcing brief and quantity")
	}
	req.BuyerName = user.Name
	req.BuyerCompany = user.CompanyName
	req.BuyerCountry = user.Country
	rfqID := fmt.Sprintf("rfq_%d", time.Now().UnixNano())
	rfqNumber := fmt.Sprintf("RFQ-%d-%d", time.Now().Year(), time.Now().Unix()%100000)

	buyerID := "guest"
	buyerName := req.BuyerName
	buyerCompany := req.BuyerCompany
	buyerCountry := req.BuyerCountry

	if user != nil {
		buyerID = user.ID
		if buyerName == "" {
			buyerName = user.Name
		}
		if buyerCompany == "" {
			buyerCompany = user.CompanyName
		}
		if buyerCountry == "" {
			buyerCountry = user.Country
		}
	}

	rfq := &models.Rfq{
		ID:                 rfqID,
		RfqNumber:          rfqNumber,
		BuyerID:            buyerID,
		BuyerName:          buyerName,
		BuyerCompany:       buyerCompany,
		BuyerCountry:       buyerCountry,
		Title:              req.Title,
		Category:           req.Category,
		Quantity:           req.Quantity,
		TargetPrice:        req.TargetPrice,
		DestinationCountry: req.DestinationCountry,
		Specifications:     req.Specifications,
		Status:             models.RfqStatusOpen,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := facades.Orm().Query().Create(rfq); err != nil {
		return nil, err
	}

	return rfq, nil
}
