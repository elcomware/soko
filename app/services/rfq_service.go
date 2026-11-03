package services

import (
	"elcomware/app/actions"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
)

type RfqService struct {
	rfqQueries        *queries.RfqQueries
	createRfqAction   *actions.CreateRfqAction
	submitQuoteAction *actions.SubmitRfqQuoteAction
}

func NewRfqService() *RfqService {
	return &RfqService{
		rfqQueries:        queries.NewRfqQueries(),
		createRfqAction:   actions.NewCreateRfqAction(),
		submitQuoteAction: actions.NewSubmitRfqQuoteAction(),
	}
}

func (s *RfqService) GetOpenRfqs() ([]models.Rfq, error) {
	return s.rfqQueries.FindAll()
}

func (s *RfqService) GetRfqByID(id string) (*models.Rfq, error) {
	return s.rfqQueries.FindByID(id)
}

func (s *RfqService) CreateRfq(req *requests.CreateRfqRequest, user *models.User) (*models.Rfq, error) {
	return s.createRfqAction.ExecuteFromRequest(req, user)
}

func (s *RfqService) SubmitQuote(req *requests.SubmitRfqQuoteRequest, user *models.User) (*models.RfqQuote, error) {
	return s.submitQuoteAction.ExecuteFromRequest(req, user)
}
