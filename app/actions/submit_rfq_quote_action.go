package actions

import (
	"fmt"
	"time"

	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
)

type SubmitRfqQuoteAction struct{}

func NewSubmitRfqQuoteAction() *SubmitRfqQuoteAction {
	return &SubmitRfqQuoteAction{}
}

func (a *SubmitRfqQuoteAction) ExecuteFromRequest(req *requests.SubmitRfqQuoteRequest, user *models.User) (*models.RfqQuote, error) {
	var rfq models.Rfq
	if err := facades.Orm().Query().Where("id = ?", req.RfqID).First(&rfq); err != nil {
		return nil, fmt.Errorf("rfq not found: %w", err)
	}

	vendorID := "usr_vendor_01"
	vendorName := "Apex Technologies"
	vendorCountry := "Germany"

	if user != nil {
		vendorID = user.ID
		if user.CompanyName != "" {
			vendorName = user.CompanyName
		} else {
			vendorName = user.Name
		}
		if user.Country != "" {
			vendorCountry = user.Country
		}
	}

	leadTime := req.LeadTimeDays
	if leadTime <= 0 {
		leadTime = 7
	}

	quote := &models.RfqQuote{
		ID:            fmt.Sprintf("quote_%d", time.Now().UnixNano()),
		RfqID:         req.RfqID,
		VendorID:      vendorID,
		VendorName:    vendorName,
		VendorCountry: vendorCountry,
		UnitPrice:     req.UnitPrice,
		TotalAmount:   req.UnitPrice * float64(rfq.Quantity),
		LeadTimeDays:  leadTime,
		Notes:         req.Notes,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := facades.Orm().Query().Create(quote); err != nil {
		return nil, err
	}

	return quote, nil
}
