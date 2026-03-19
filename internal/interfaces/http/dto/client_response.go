package dto

type ClientResponse struct {
	ID               uint    `json:"id"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	Active           bool    `json:"active"`
	AppointmentCount int64   `json:"appointment_count"`
	TotalSpent       float64 `json:"total_spent"`
}
