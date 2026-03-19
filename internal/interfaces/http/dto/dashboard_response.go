package dto

import "apiGoShei/internal/domain/appointment"

type PeriodStatsResponse struct {
	Total     int64   `json:"total"`
	Completed int64   `json:"completed"`
	Cancelled int64   `json:"cancelled"`
	Pending   int64   `json:"pending"`
	Confirmed int64   `json:"confirmed"`
	Revenue   float64 `json:"revenue"`
	Deposits  float64 `json:"deposits"`
}

type MetricsResponse struct {
	Today PeriodStatsResponse `json:"today"`
	Week  PeriodStatsResponse `json:"week"`
	Month PeriodStatsResponse `json:"month"`
	Year  PeriodStatsResponse `json:"year"`
}

func toPeriodResponse(p appointment.PeriodStats) PeriodStatsResponse {
	return PeriodStatsResponse{
		Total:     p.Total,
		Completed: p.Completed,
		Cancelled: p.Cancelled,
		Pending:   p.Pending,
		Confirmed: p.Confirmed,
		Revenue:   p.Revenue,
		Deposits:  p.Deposits,
	}
}

func MetricsToResponse(m *appointment.DashboardMetrics) MetricsResponse {
	return MetricsResponse{
		Today: toPeriodResponse(m.Today),
		Week:  toPeriodResponse(m.Week),
		Month: toPeriodResponse(m.Month),
		Year:  toPeriodResponse(m.Year),
	}
}
