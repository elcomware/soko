package services

import (
	"elcomware/app/actions"
	"elcomware/app/dto"
	"elcomware/app/mappers"
	"elcomware/app/queries"
)

type OfferService struct {
	offerQueries      *queries.OfferQueries
	offerMapper       *mappers.OfferMapper
	submitOfferAction *actions.SubmitOfferAction
}

func NewOfferService() *OfferService {
	return &OfferService{
		offerQueries:      queries.NewOfferQueries(),
		offerMapper:       mappers.NewOfferMapper(),
		submitOfferAction: actions.NewSubmitOfferAction(),
	}
}

func (s *OfferService) GetPendingOffers() ([]dto.OfferOutput, error) {
	offers, err := s.offerQueries.FindAll()
	if err != nil {
		return nil, err
	}
	return s.offerMapper.ToDTOList(offers), nil
}

func (s *OfferService) GetOffersByProductID(productID string) ([]dto.OfferOutput, error) {
	offers, err := s.offerQueries.FindByProductID(productID)
	if err != nil {
		return nil, err
	}
	return s.offerMapper.ToDTOList(offers), nil
}

func (s *OfferService) SubmitOffer(input dto.SubmitOfferInput) (dto.OfferOutput, error) {
	return s.submitOfferAction.Execute(input)
}
