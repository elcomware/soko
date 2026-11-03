package services

import (
	"elcomware/app/actions"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
)

type OfferService struct {
	offerQueries      *queries.OfferQueries
	submitOfferAction *actions.SubmitOfferAction
}

func NewOfferService() *OfferService {
	return &OfferService{
		offerQueries:      queries.NewOfferQueries(),
		submitOfferAction: actions.NewSubmitOfferAction(),
	}
}

func (s *OfferService) GetOffers() ([]models.Offer, error) {
	return s.offerQueries.FindAll()
}

func (s *OfferService) GetOffersByProductID(productID string) ([]models.Offer, error) {
	return s.offerQueries.FindByProductID(productID)
}

func (s *OfferService) SubmitOffer(req *requests.SubmitOfferRequest, user *models.User) (*models.Offer, error) {
	return s.submitOfferAction.ExecuteFromRequest(req, user)
}
