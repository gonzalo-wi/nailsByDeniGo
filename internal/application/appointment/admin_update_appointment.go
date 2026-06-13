package appointmentapp

import (
	"errors"

	"apiGoShei/internal/domain/appointment"
	"apiGoShei/internal/domain/service"
)

type AdminUpdateAppointmentInput struct {
	ID            uint
	Status        *appointment.AppointmentStatus
	ServiceID     *uint
	PenaltyAmount *float64
	PenaltyNote   string
}

var ErrInvalidStatus = errors.New("estado inválido")

type AdminUpdateAppointmentUseCase struct {
	appointmentRepo appointment.Repository
	serviceRepo     service.Repository
}

func NewAdminUpdateAppointmentUseCase(apptRepo appointment.Repository, svcRepo service.Repository) *AdminUpdateAppointmentUseCase {
	return &AdminUpdateAppointmentUseCase{appointmentRepo: apptRepo, serviceRepo: svcRepo}
}

func (uc *AdminUpdateAppointmentUseCase) Execute(input AdminUpdateAppointmentInput) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}

	if input.Status != nil {
		switch *input.Status {
		case appointment.StatusPending, appointment.StatusConfirmed, appointment.StatusDone,
			appointment.StatusCancelled, appointment.StatusAbsent:
		default:
			return nil, ErrInvalidStatus
		}
		appt.Status = *input.Status
	}

	if input.ServiceID != nil && *input.ServiceID != appt.ServiceID {
		svc, err := uc.serviceRepo.FindByID(*input.ServiceID)
		if err != nil {
			return nil, err
		}
		if svc == nil {
			return nil, service.ErrNotFound
		}
		appt.ServiceID = svc.ID
		appt.BasePrice = svc.BasePrice
		appt.FinalPrice = svc.BasePrice + appt.ExtrasAmount
	}

	if input.PenaltyAmount != nil {
		appt.PenaltyAmount = *input.PenaltyAmount
		appt.PenaltyNote = input.PenaltyNote
	}

	if err := uc.appointmentRepo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}
