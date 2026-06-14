package appointmentapp

import (
	"time"

	"apiGoShei/internal/domain/appointment"
)

type ListAppointmentUseCase struct {
	appointmentRepo appointment.Repository
}

func NewListAppointmentUseCase(repo appointment.Repository) *ListAppointmentUseCase {
	return &ListAppointmentUseCase{appointmentRepo: repo}
}

func (uc *ListAppointmentUseCase) Execute(filters appointment.AppointmentFilters) ([]appointment.Appointment, error) {
	return uc.appointmentRepo.FindAll(filters)
}

type NextAppointmentUseCase struct {
	appointmentRepo appointment.Repository
}

func NewNextAppointmentUseCase(repo appointment.Repository) *NextAppointmentUseCase {
	return &NextAppointmentUseCase{appointmentRepo: repo}
}

func (uc *NextAppointmentUseCase) Execute(clientID uint) (*appointment.Appointment, error) {
	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	today := time.Now().In(loc).Format("2006-01-02")
	return uc.appointmentRepo.FindNextByClient(clientID, today)
}

type GetAppointmentByIDUseCase struct {
	appointmentRepo appointment.Repository
}

func NewGetAppointmentByIDUseCase(repo appointment.Repository) *GetAppointmentByIDUseCase {
	return &GetAppointmentByIDUseCase{appointmentRepo: repo}
}

func (uc *GetAppointmentByIDUseCase) Execute(id uint) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}
	return appt, nil
}

type CompleteAppointmentUseCase struct {
	appointmentRepo appointment.Repository
}

func NewCompleteAppointmentUseCase(repo appointment.Repository) *CompleteAppointmentUseCase {
	return &CompleteAppointmentUseCase{appointmentRepo: repo}
}

func (uc *CompleteAppointmentUseCase) Execute(id uint) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}
	if appt.Status != appointment.StatusConfirmed && appt.Status != appointment.StatusPending {
		return nil, appointment.ErrInvalidStatus
	}
	appt.Status = appointment.StatusDone
	if err := uc.appointmentRepo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}
