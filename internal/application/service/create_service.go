package serviceapp

import "apiGoShei/internal/domain/service"

type CreateServiceInput struct {
	Name             string
	Description      string
	DurationMinutes  int
	BasePrice        float64
	RequiresDeposit  bool
	SuggestedDeposit float64
	Color            string
}

type CreateServiceUseCase struct {
	serviceRepo service.Repository
}

func NewCreateServiceUseCase(repo service.Repository) *CreateServiceUseCase {
	return &CreateServiceUseCase{serviceRepo: repo}
}

func (uc *CreateServiceUseCase) Execute(input CreateServiceInput) (*service.Service, error) {
	s := &service.Service{
		Name:             input.Name,
		Description:      input.Description,
		DurationMinutes:  input.DurationMinutes,
		BasePrice:        input.BasePrice,
		RequiresDeposit:  input.RequiresDeposit,
		SuggestedDeposit: input.SuggestedDeposit,
		Color:            input.Color,
		Active:           true,
	}
	if err := uc.serviceRepo.Create(s); err != nil {
		return nil, err
	}
	return s, nil
}
