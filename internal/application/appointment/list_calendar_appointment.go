package appointmentapp

import "apiGoShei/internal/domain/appointment"

type ListCalendarAppointmentInput struct {
	From string
	To   string
}

type ListCalendarAppointmentUseCase struct {
	appointmentRepo appointment.Repository
}

func NewListCalendarAppointmentUseCase(repo appointment.Repository) *ListCalendarAppointmentUseCase {
	return &ListCalendarAppointmentUseCase{appointmentRepo: repo}
}

func (uc *ListCalendarAppointmentUseCase) Execute(input ListCalendarAppointmentInput) ([]appointment.Appointment, error) {
	return uc.appointmentRepo.FindByDateRange(input.From, input.To)
}
