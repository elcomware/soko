package dto

type DashboardStatsOutput struct {
	TotalRevenue        float64 `json:"total_revenue"`
	ActiveOrders        int64   `json:"active_orders"`
	TotalProducts       int64   `json:"total_products"`
	PendingRfqs         int64   `json:"pending_rfqs"`
	ActiveNegotiations  int64   `json:"active_negotiations"`
	TradeVolumeCovered  float64 `json:"trade_volume_covered"`
}
