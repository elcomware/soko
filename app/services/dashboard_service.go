package services

import (
	"elcomware/app/dto"
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

func (s *DashboardService) GetStats() (dto.DashboardStatsOutput, error) {
	return s.dashboardQueries.GetStats()
}
