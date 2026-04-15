package appointmentapp

import "apiGoShei/internal/domain/appointment"

type UpdateFinalPriceInput struct {
	ID           uint
	ExtrasAmount float64
	ExtrasNote   string
}

type UpdateFinalPriceUseCase struct {
	appointmentRepo appointment.Repository
}

func NewUpdateFinalPriceUseCase(repo appointment.Repository) *UpdateFinalPriceUseCase {
	return &UpdateFinalPriceUseCase{appointmentRepo: repo}
}

func (uc *UpdateFinalPriceUseCase) Execute(input UpdateFinalPriceInput) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}
	appt.ExtrasAmount = input.ExtrasAmount
	appt.ExtrasNote = input.ExtrasNote
	appt.FinalPrice = appt.BasePrice + input.ExtrasAmount
	if err := uc.appointmentRepo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}

type UpdateDepositInput struct {
	ID            uint
	DepositAmount float64
}

type UpdateDepositUseCase struct {
	appointmentRepo appointment.Repository
}

func NewUpdateDepositUseCase(repo appointment.Repository) *UpdateDepositUseCase {
	return &UpdateDepositUseCase{appointmentRepo: repo}
}

func (uc *UpdateDepositUseCase) Execute(input UpdateDepositInput) (*appointment.Appointment, error) {
	appt, err := uc.appointmentRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, appointment.ErrNotFound
	}
	appt.DepositAmount = input.DepositAmount
	if err := uc.appointmentRepo.Update(appt); err != nil {
		return nil, err
	}
	return appt, nil
}
