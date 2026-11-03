package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
)

type RfqQueries struct{}

func NewRfqQueries() *RfqQueries {
	return &RfqQueries{}
}

func (q *RfqQueries) FindAll() ([]models.Rfq, error) {
	var rfqs []models.Rfq
	if err := facades.Orm().Query().OrderByDesc("created_at").Get(&rfqs); err != nil {
		return nil, err
	}

	for i := range rfqs {
		var quotes []models.RfqQuote
		_ = facades.Orm().Query().Where("rfq_id = ?", rfqs[i].ID).OrderByDesc("created_at").Get(&quotes)
		rfqs[i].Quotes = quotes
	}

	return rfqs, nil
}

func (q *RfqQueries) FindByID(id string) (*models.Rfq, error) {
	var rfq models.Rfq
	if err := facades.Orm().Query().Where("id = ?", id).First(&rfq); err != nil {
		return nil, err
	}

	var quotes []models.RfqQuote
	_ = facades.Orm().Query().Where("rfq_id = ?", rfq.ID).OrderByDesc("created_at").Get(&quotes)
	rfq.Quotes = quotes

	return &rfq, nil
}

func (q *RfqQueries) CountPending() (int64, error) {
	return facades.Orm().Query().Model(&models.Rfq{}).Where("status = ?", models.RfqStatusOpen).Count()
}
