package dto

import "apiGoShei/internal/domain/service"

type ServiceResponse struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	DurationMinutes  int     `json:"duration_minutes"`
	BasePrice        float64 `json:"base_price"`
	RequiresDeposit  bool    `json:"requires_deposit"`
	SuggestedDeposit float64 `json:"suggested_deposit"`
	Color            string  `json:"color"`
	Active           bool    `json:"active"`
}

func ServiceToResponse(s *service.Service) ServiceResponse {
	return ServiceResponse{
		ID:               s.ID,
		Name:             s.Name,
		Description:      s.Description,
		DurationMinutes:  s.DurationMinutes,
		BasePrice:        s.BasePrice,
		RequiresDeposit:  s.RequiresDeposit,
		SuggestedDeposit: s.SuggestedDeposit,
		Color:            s.Color,
		Active:           s.Active,
	}
}
