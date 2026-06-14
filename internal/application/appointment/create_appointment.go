package appointmentapp

import (
	"errors"
	"time"

	"apiGoShei/internal/domain/appointment"
	"apiGoShei/internal/domain/client"
	"apiGoShei/internal/domain/schedule"
	"apiGoShei/internal/domain/service"
	"apiGoShei/internal/infrastructure/logger"
	"apiGoShei/internal/infrastructure/mail"
)

type CreateAppointmentInput struct {
	ClientID       uint
	ServiceID      uint
	ProfessionalID *uint
	Date           string
	StartTime      string
	Notes          string
	IsAdmin        bool
}

type CreateAppointmentUseCase struct {
	appointmentRepo appointment.Repository
	clientRepo      client.Repository
	serviceRepo     service.Repository
	scheduleRepo    schedule.Repository
	mailer          mail.Mailer
}

func NewCreateAppointmentUseCase(
	appointmentRepo appointment.Repository,
	clientRepo client.Repository,
	serviceRepo service.Repository,
	scheduleRepo schedule.Repository,
	mailer mail.Mailer,
) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{
		appointmentRepo: appointmentRepo,
		clientRepo:      clientRepo,
		serviceRepo:     serviceRepo,
		scheduleRepo:    scheduleRepo,
		mailer:          mailer,
	}
}

func (uc *CreateAppointmentUseCase) Execute(input CreateAppointmentInput) (*appointment.Appointment, error) {
	if !input.IsAdmin {
		loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
		today := time.Now().In(loc).Format("2006-01-02")
		if input.Date < today {
			return nil, appointment.ErrPastDate
		}
	}

	clientFound, err := uc.clientRepo.FindByID(input.ClientID)
	if err != nil {
		return nil, err
	}
	if clientFound == nil {
		return nil, errors.New("clienta no encontrada")
	}

	serviceFound, err := uc.serviceRepo.FindByID(input.ServiceID)
	if err != nil {
		return nil, err
	}
	if serviceFound == nil {
		return nil, errors.New("servicio no encontrado")
	}

	hasAppt, err := uc.appointmentRepo.ExistsByClientAndDate(input.ClientID, input.Date)
	if err != nil {
		return nil, err
	}
	if hasAppt {
		return nil, appointment.ErrClientHasThatDay
	}

	available, err := uc.scheduleRepo.IsAvailable(input.Date, input.StartTime, serviceFound.DurationMinutes)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, appointment.ErrSlotUnavailable
	}

	endTime := calculateEndTime(input.StartTime, serviceFound.DurationMinutes)

	overlap, err := uc.appointmentRepo.ExistsOverlap(input.Date, input.StartTime, endTime, input.ProfessionalID)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, appointment.ErrOverlap
	}

	newAppointment := &appointment.Appointment{
		ClientID:       input.ClientID,
		ServiceID:      input.ServiceID,
		ProfessionalID: input.ProfessionalID,
		Date:           input.Date,
		StartTime:      input.StartTime,
		EndTime:        endTime,
		BasePrice:      serviceFound.BasePrice,
		ExtrasAmount:   0,
		DepositAmount:  serviceFound.SuggestedDeposit,
		FinalPrice:     serviceFound.BasePrice,
		Status:         appointment.StatusPending,
		Notes:          input.Notes,
	}

	if err := uc.appointmentRepo.Create(newAppointment); err != nil {
		return nil, err
	}

	// Enviar mails de forma asíncrona para no bloquear la respuesta HTTP.
	go uc.sendNotifications(clientFound.FirstName+" "+clientFound.LastName, clientFound.Email, clientFound.Phone, serviceFound.Name, newAppointment)

	return newAppointment, nil
}

func (uc *CreateAppointmentUseCase) sendNotifications(clientName, clientEmail, clientPhone, serviceName string, appt *appointment.Appointment) {
	data := mail.AppointmentData{
		ClientName:    clientName,
		ClientEmail:   clientEmail,
		ClientPhone:   clientPhone,
		ServiceName:   serviceName,
		Date:          appt.Date,
		StartTime:     appt.StartTime,
		EndTime:       appt.EndTime,
		Status:        string(appt.Status),
		Notes:         appt.Notes,
		BasePrice:     appt.BasePrice,
		ExtrasAmount:  appt.ExtrasAmount,
		ExtrasNote:    appt.ExtrasNote,
		FinalPrice:    appt.FinalPrice,
		DepositAmount: appt.DepositAmount,
	}
	if err := uc.mailer.SendNewAppointmentToClient(data); err != nil {
		logger.Error.Printf("error enviando mail a cliente: %v", err)
	}
	if err := uc.mailer.SendNewAppointmentToAdmin(data); err != nil {
		logger.Error.Printf("error enviando mail a admin: %v", err)
	}
}
