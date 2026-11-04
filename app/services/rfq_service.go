package services

import (
	"elcomware/app/actions"
	"elcomware/app/facades"
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
	var rfqs []models.Rfq
	if err := facades.Orm().Query().Where("status = ?", models.RfqStatusOpen).OrderByDesc("created_at").Get(&rfqs); err != nil {
		return nil, err
	}
	for i := range rfqs {
		rfqs[i].BuyerID = ""
		rfqs[i].BuyerName = ""
		rfqs[i].BuyerCompany = ""
		rfqs[i].Quotes = []models.RfqQuote{}
	}
	return rfqs, nil
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
