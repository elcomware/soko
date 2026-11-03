package services

import (
	"elcomware/app/queries"
)

type DashboardService struct {
	dashboardQueries *queries.DashboardQueries
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		dashboardQueries: queries.NewDashboardQueries(),
	}
}

func (s *DashboardService) GetStats() (map[string]any, error) {
	return s.dashboardQueries.GetStats()
}
