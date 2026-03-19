package appointmentapp

import "apiGoShei/internal/domain/appointment"

type ConfirmAppointmentUseCase struct {
	appointmentRepo appointment.Repository
}

func NewConfirmAppointmentUseCase(repo appointment.Repository) *ConfirmAppointmentUseCase {
	return &ConfirmAppointmentUseCase{appointmentRepo: repo}
}

func (uc *ConfirmAppointmentUseCase) Execute(id uint) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}
	if appt.Status != appointment.StatusPending {
		return nil, appointment.ErrInvalidStatus
	}
	appt.Status = appointment.StatusConfirmed
	if err := uc.appointmentRepo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}
