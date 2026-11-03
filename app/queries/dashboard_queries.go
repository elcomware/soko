package queries

import (
	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/models"
)

type DashboardQueries struct {
	productQueries *ProductQueries
	orderQueries   *OrderQueries
	rfqQueries     *RfqQueries
	offerQueries   *OfferQueries
}

func NewDashboardQueries() *DashboardQueries {
	return &DashboardQueries{
		productQueries: NewProductQueries(),
		orderQueries:   NewOrderQueries(),
		rfqQueries:     NewRfqQueries(),
		offerQueries:   NewOfferQueries(),
	}
}

func (q *DashboardQueries) GetStats() (dto.DashboardStatsOutput, error) {
	revenue, _ := q.orderQueries.SumRevenue()
	activeOrders, _ := q.orderQueries.CountActive()
	totalProducts, _ := q.productQueries.Count()
	pendingRfqs, _ := q.rfqQueries.CountPending()
	activeOffers, _ := q.offerQueries.CountActive()

	// Trade Assurance Volume covered across all products
	var tradeCoverTotal float64
	_ = facades.Orm().Query().Model(&models.Product{}).Sum("trade_assurance_cover", &tradeCoverTotal)

	return dto.DashboardStatsOutput{
		TotalRevenue:       revenue,
		ActiveOrders:       activeOrders,
		TotalProducts:      totalProducts,
		PendingRfqs:        pendingRfqs,
		ActiveNegotiations: activeOffers,
		TradeVolumeCovered: tradeCoverTotal,
	}, nil
}
