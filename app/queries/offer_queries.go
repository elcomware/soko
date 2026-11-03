package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
)

type OfferQueries struct{}

func NewOfferQueries() *OfferQueries {
	return &OfferQueries{}
}

func (q *OfferQueries) FindAll() ([]models.Offer, error) {
	var offers []models.Offer
	err := facades.Orm().Query().OrderByDesc("created_at").Get(&offers)
	return offers, err
}

func (q *OfferQueries) FindByProductID(productID string) ([]models.Offer, error) {
	var offers []models.Offer
	err := facades.Orm().Query().Where("product_id = ?", productID).OrderByDesc("created_at").Get(&offers)
	return offers, err
}

func (q *OfferQueries) CountActive() (int64, error) {
	return facades.Orm().Query().Model(&models.Offer{}).Where("status = ?", models.OfferStatusPending).Count()
}
