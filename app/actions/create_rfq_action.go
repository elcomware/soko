package actions

import (
	"fmt"
	"time"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type CreateRfqAction struct {
	rfqMapper *mappers.RfqMapper
}

func NewCreateRfqAction() *CreateRfqAction {
	return &CreateRfqAction{
		rfqMapper: mappers.NewRfqMapper(),
	}
}

func (a *CreateRfqAction) Execute(input dto.CreateRfqInput) (dto.RfqOutput, error) {
	rfqID := fmt.Sprintf("rfq_%d", time.Now().UnixNano())
	rfqNumber := fmt.Sprintf("RFQ-%d-%d", time.Now().Year(), time.Now().Unix()%100000)

	rfq := models.Rfq{
		ID:                 rfqID,
		RfqNumber:          rfqNumber,
		BuyerID:            input.BuyerID,
		BuyerName:          input.BuyerName,
		BuyerCompany:       input.BuyerCompany,
		BuyerCountry:       input.BuyerCountry,
		Title:              input.Title,
		Category:           input.Category,
		Quantity:           input.Quantity,
		TargetPrice:        input.TargetPrice,
		DestinationCountry: input.DestinationCountry,
		Specifications:     input.Specifications,
		Status:             models.RfqStatusOpen,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := facades.Orm().Query().Create(&rfq); err != nil {
		return dto.RfqOutput{}, err
	}

	return a.rfqMapper.ToDTO(rfq), nil
}
