package dto

type CreateServiceRequest struct {
	Name             string  `json:"name"              binding:"required"`
	Description      string  `json:"description"`
	DurationMinutes  int     `json:"duration_minutes"  binding:"required,min=1"`
	BasePrice        float64 `json:"base_price"        binding:"required,min=0"`
	RequiresDeposit  bool    `json:"requires_deposit"`
	SuggestedDeposit float64 `json:"suggested_deposit"`
	Color            string  `json:"color"`
}

type UpdateServiceRequest struct {
	Name             string  `json:"name"              binding:"required"`
	Description      string  `json:"description"`
	DurationMinutes  int     `json:"duration_minutes"  binding:"required,min=1"`
	BasePrice        float64 `json:"base_price"        binding:"required,min=0"`
	RequiresDeposit  bool    `json:"requires_deposit"`
	SuggestedDeposit float64 `json:"suggested_deposit"`
	Color            string  `json:"color"`
}
