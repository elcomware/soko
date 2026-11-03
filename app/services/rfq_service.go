package services

import (
	"elcomware/app/actions"
	"elcomware/app/dto"
	"elcomware/app/mappers"
	"elcomware/app/queries"
)

type RfqService struct {
	rfqQueries      *queries.RfqQueries
	rfqMapper       *mappers.RfqMapper
	createRfqAction *actions.CreateRfqAction
	submitQuoteAction *actions.SubmitRfqQuoteAction
}

func NewRfqService() *RfqService {
	return &RfqService{
		rfqQueries:      queries.NewRfqQueries(),
		rfqMapper:       mappers.NewRfqMapper(),
		createRfqAction: actions.NewCreateRfqAction(),
		submitQuoteAction: actions.NewSubmitRfqQuoteAction(),
	}
}

func (s *RfqService) GetOpenRfqs() ([]dto.RfqOutput, error) {
	rfqs, err := s.rfqQueries.FindAll()
	if err != nil {
		return nil, err
	}
	return s.rfqMapper.ToDTOList(rfqs), nil
}

func (s *RfqService) GetRfqByID(id string) (*dto.RfqOutput, error) {
	rfq, err := s.rfqQueries.FindByID(id)
	if err != nil {
		return nil, err
	}
	dto := s.rfqMapper.ToDTO(*rfq)
	return &dto, nil
}

func (s *RfqService) CreateRfq(input dto.CreateRfqInput) (dto.RfqOutput, error) {
	return s.createRfqAction.Execute(input)
}

func (s *RfqService) SubmitQuote(input dto.SubmitQuoteInput) (dto.RfqQuoteOutput, error) {
	return s.submitQuoteAction.Execute(input)
}
