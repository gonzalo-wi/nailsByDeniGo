package appointment

import "time"

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "PENDING"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusDone      AppointmentStatus = "DONE"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusAbsent    AppointmentStatus = "ABSENT"
)

type Appointment struct {
	ID             uint
	ClientID       uint
	ServiceID      uint
	ProfessionalID *uint
	Date           string
	StartTime      string
	EndTime        string
	BasePrice      float64
	ExtrasAmount   float64
	ExtrasNote     string
	DepositAmount  float64
	FinalPrice     float64
	Status         AppointmentStatus
	Notes          string
	PenaltyAmount  float64
	PenaltyNote    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
