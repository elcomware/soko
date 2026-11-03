package mappers

import (
	"elcomware/app/dto"
	"elcomware/app/models"
)

type RfqMapper struct{}

func NewRfqMapper() *RfqMapper {
	return &RfqMapper{}
}

func (m *RfqMapper) ToQuoteDTO(q models.RfqQuote) dto.RfqQuoteOutput {
	return dto.RfqQuoteOutput{
		ID:            q.ID,
		RfqID:         q.RfqID,
		VendorID:      q.VendorID,
		VendorName:    q.VendorName,
		VendorCountry: q.VendorCountry,
		UnitPrice:     q.UnitPrice,
		TotalAmount:   q.TotalAmount,
		LeadTimeDays:  q.LeadTimeDays,
		Notes:         q.Notes,
		Status:        q.Status,
		CreatedAt:     q.CreatedAt,
	}
}

func (m *RfqMapper) ToDTO(rfq models.Rfq) dto.RfqOutput {
	quotes := make([]dto.RfqQuoteOutput, 0, len(rfq.Quotes))
	for _, q := range rfq.Quotes {
		quotes = append(quotes, m.ToQuoteDTO(q))
	}

	return dto.RfqOutput{
		ID:                 rfq.ID,
		RfqNumber:          rfq.RfqNumber,
		BuyerID:            rfq.BuyerID,
		BuyerName:          rfq.BuyerName,
		BuyerCompany:       rfq.BuyerCompany,
		BuyerCountry:       rfq.BuyerCountry,
		Title:              rfq.Title,
		Category:           rfq.Category,
		Quantity:           rfq.Quantity,
		TargetPrice:        rfq.TargetPrice,
		DestinationCountry: rfq.DestinationCountry,
		Specifications:     rfq.Specifications,
		Status:             string(rfq.Status),
		Quotes:             quotes,
		CreatedAt:          rfq.CreatedAt,
	}
}

func (m *RfqMapper) ToDTOList(list []models.Rfq) []dto.RfqOutput {
	dtos := make([]dto.RfqOutput, 0, len(list))
	for _, r := range list {
		dtos = append(dtos, m.ToDTO(r))
	}
	return dtos
}
