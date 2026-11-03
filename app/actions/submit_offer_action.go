package actions

import (
	"fmt"
	"time"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type SubmitOfferAction struct {
	offerMapper *mappers.OfferMapper
}

func NewSubmitOfferAction() *SubmitOfferAction {
	return &SubmitOfferAction{
		offerMapper: mappers.NewOfferMapper(),
	}
}

func (a *SubmitOfferAction) Execute(input dto.SubmitOfferInput) (dto.OfferOutput, error) {
	var product models.Product
	if err := facades.Orm().Query().Where("id = ?", input.ProductID).First(&product); err != nil {
		return dto.OfferOutput{}, fmt.Errorf("product %s not found: %w", input.ProductID, err)
	}

	totalOffered := input.OfferedPrice * float64(input.Quantity)
	offer := models.Offer{
		ID:           fmt.Sprintf("off_%d", time.Now().UnixNano()),
		ProductID:    product.ID,
		ProductTitle: product.Title,
		BuyerID:      input.BuyerID,
		BuyerName:    input.BuyerName,
		BuyerEmail:   input.BuyerEmail,
		BuyerCompany: input.BuyerCompany,
		OfferedPrice: input.OfferedPrice,
		Quantity:     input.Quantity,
		TotalOffered: totalOffered,
		Message:      input.Message,
		Status:       models.OfferStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := facades.Orm().Query().Create(&offer); err != nil {
		return dto.OfferOutput{}, err
	}

	return a.offerMapper.ToDTO(offer), nil
}
