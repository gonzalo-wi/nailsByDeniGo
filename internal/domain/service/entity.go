package service

import "time"

type Service struct {
	ID               uint
	Name             string
	Description      string
	DurationMinutes  int
	BasePrice        float64
	RequiresDeposit  bool
	SuggestedDeposit float64
	Color            string
	Active           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
