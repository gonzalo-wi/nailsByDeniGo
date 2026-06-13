package dto

import "apiGoShei/internal/domain/appointment"

type AppointmentResponse struct {
	ID             uint    `json:"id"`
	ClientID       uint    `json:"client_id"`
	ServiceID      uint    `json:"service_id"`
	ProfessionalID *uint   `json:"professional_id,omitempty"`
	Date           string  `json:"date"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	BasePrice      float64 `json:"base_price"`
	ExtrasAmount   float64 `json:"extras_amount"`
	ExtrasNote     string  `json:"extras_note,omitempty"`
	DepositAmount  float64 `json:"deposit_amount"`
	FinalPrice     float64 `json:"final_price"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes,omitempty"`
	PenaltyAmount  float64 `json:"penalty_amount"`
	PenaltyNote    string  `json:"penalty_note,omitempty"`
}

func AppointmentToResponse(a *appointment.Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:             a.ID,
		ClientID:       a.ClientID,
		ServiceID:      a.ServiceID,
		ProfessionalID: a.ProfessionalID,
		Date:           a.Date,
		StartTime:      a.StartTime,
		EndTime:        a.EndTime,
		BasePrice:      a.BasePrice,
		ExtrasAmount:   a.ExtrasAmount,
		ExtrasNote:     a.ExtrasNote,
		DepositAmount:  a.DepositAmount,
		FinalPrice:     a.FinalPrice,
		Status:         string(a.Status),
		Notes:          a.Notes,
		PenaltyAmount:  a.PenaltyAmount,
		PenaltyNote:    a.PenaltyNote,
	}
}
