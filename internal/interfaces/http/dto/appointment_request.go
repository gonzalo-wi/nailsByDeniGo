package dto

type CreateAppointmentRequest struct {
	ClientID       uint   `json:"client_id"        binding:"required"`
	ServiceID      uint   `json:"service_id"       binding:"required"`
	ProfessionalID *uint  `json:"professional_id"`
	Date           string `json:"date"             binding:"required"`
	StartTime      string `json:"start_time"       binding:"required"`
	Notes          string `json:"notes"`
}

type UpdateFinalPriceRequest struct {
	ExtrasAmount float64 `json:"extras_amount" binding:"required,min=0"`
	ExtrasNote   string  `json:"extras_note"`
}

type UpdateDepositRequest struct {
	DepositAmount float64 `json:"deposit_amount" binding:"required,min=0"`
}
