package appointmentapp

import "apiGoShei/internal/domain/appointment"

type CancelAppointmentUseCase struct {
	appointmentRepo appointment.Repository
}

func NewCancelAppointmentUseCase(repo appointment.Repository) *CancelAppointmentUseCase {
	return &CancelAppointmentUseCase{appointmentRepo: repo}
}

func (uc *CancelAppointmentUseCase) Execute(id uint) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}
	if appt.Status == appointment.StatusDone || appt.Status == appointment.StatusCancelled {
		return nil, appointment.ErrInvalidStatus
	}
	appt.Status = appointment.StatusCancelled
	if err := uc.appointmentRepo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}
