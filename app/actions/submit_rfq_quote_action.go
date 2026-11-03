package actions

import (
	"fmt"
	"time"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type SubmitRfqQuoteAction struct {
	rfqMapper *mappers.RfqMapper
}

func NewSubmitRfqQuoteAction() *SubmitRfqQuoteAction {
	return &SubmitRfqQuoteAction{
		rfqMapper: mappers.NewRfqMapper(),
	}
}

func (a *SubmitRfqQuoteAction) Execute(input dto.SubmitQuoteInput) (dto.RfqQuoteOutput, error) {
	var rfq models.Rfq
	if err := facades.Orm().Query().Where("id = ?", input.RfqID).First(&rfq); err != nil {
		return dto.RfqQuoteOutput{}, fmt.Errorf("rfq with ID %s not found: %w", input.RfqID, err)
	}

	totalAmount := input.UnitPrice * float64(rfq.Quantity)
	quote := models.RfqQuote{
		ID:            fmt.Sprintf("quote_%d", time.Now().UnixNano()),
		RfqID:         input.RfqID,
		VendorID:      input.VendorID,
		VendorName:    input.VendorName,
		VendorCountry: input.VendorCountry,
		UnitPrice:     input.UnitPrice,
		TotalAmount:   totalAmount,
		LeadTimeDays:  input.LeadTimeDays,
		Notes:         input.Notes,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := facades.Orm().Query().Create(&quote); err != nil {
		return dto.RfqQuoteOutput{}, err
	}

	return a.rfqMapper.ToQuoteDTO(quote), nil
}
