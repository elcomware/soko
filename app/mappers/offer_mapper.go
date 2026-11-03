package mappers

import (
	"elcomware/app/dto"
	"elcomware/app/models"
)

type OfferMapper struct{}

func NewOfferMapper() *OfferMapper {
	return &OfferMapper{}
}

func (m *OfferMapper) ToDTO(o models.Offer) dto.OfferOutput {
	return dto.OfferOutput{
		ID:           o.ID,
		ProductID:    o.ProductID,
		ProductTitle: o.ProductTitle,
		BuyerID:      o.BuyerID,
		BuyerName:    o.BuyerName,
		BuyerEmail:   o.BuyerEmail,
		BuyerCompany: o.BuyerCompany,
		OfferedPrice: o.OfferedPrice,
		Quantity:     o.Quantity,
		TotalOffered: o.TotalOffered,
		Message:      o.Message,
		Status:       string(o.Status),
		CreatedAt:    o.CreatedAt,
	}
}

func (m *OfferMapper) ToDTOList(list []models.Offer) []dto.OfferOutput {
	dtos := make([]dto.OfferOutput, 0, len(list))
	for _, o := range list {
		dtos = append(dtos, m.ToDTO(o))
	}
	return dtos
}
