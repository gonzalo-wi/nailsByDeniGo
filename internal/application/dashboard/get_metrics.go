package dashboardapp

import (
	"time"

	"apiGoShei/internal/domain/appointment"
)

type GetMetricsUseCase struct {
	appointmentRepo appointment.Repository
}

func NewGetMetricsUseCase(repo appointment.Repository) *GetMetricsUseCase {
	return &GetMetricsUseCase{appointmentRepo: repo}
}

func (uc *GetMetricsUseCase) Execute() (*appointment.DashboardMetrics, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1)).Format("2006-01-02")

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	return uc.appointmentRepo.GetMetrics(today, weekStart, monthStart, yearStart)
}
