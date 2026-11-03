package queries

import (
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

func (q *DashboardQueries) GetStats() (map[string]any, error) {
	revenue, _ := q.orderQueries.SumRevenue()
	activeOrders, _ := q.orderQueries.CountActive()
	totalProducts, _ := q.productQueries.Count()
	pendingRfqs, _ := q.rfqQueries.CountPending()
	activeOffers, _ := q.offerQueries.CountActive()

	var tradeCoverTotal float64
	_ = facades.Orm().Query().Model(&models.Product{}).Sum("trade_assurance_cover", &tradeCoverTotal)

	return map[string]any{
		"total_revenue":        revenue,
		"active_orders":        activeOrders,
		"total_products":       totalProducts,
		"pending_rfqs":         pendingRfqs,
		"active_negotiations":  activeOffers,
		"trade_volume_covered": tradeCoverTotal,
	}, nil
}
