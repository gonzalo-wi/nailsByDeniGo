package serviceapp

import "apiGoShei/internal/domain/service"

type UpdateServiceInput struct {
	ID               uint
	Name             string
	Description      string
	DurationMinutes  int
	BasePrice        float64
	RequiresDeposit  bool
	SuggestedDeposit float64
	Color            string
}

type UpdateServiceUseCase struct {
	serviceRepo service.Repository
}

func NewUpdateServiceUseCase(repo service.Repository) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{serviceRepo: repo}
}

func (uc *UpdateServiceUseCase) Execute(input UpdateServiceInput) (*service.Service, error) {
	s, err := uc.serviceRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, service.ErrNotFound
	}

	s.Name = input.Name
	s.Description = input.Description
	s.DurationMinutes = input.DurationMinutes
	s.BasePrice = input.BasePrice
	s.RequiresDeposit = input.RequiresDeposit
	s.SuggestedDeposit = input.SuggestedDeposit
	s.Color = input.Color

	if err := uc.serviceRepo.Update(s); err != nil {
		return nil, err
	}
	return s, nil
}

type ToggleServiceUseCase struct {
	serviceRepo service.Repository
}

func NewToggleServiceUseCase(repo service.Repository) *ToggleServiceUseCase {
	return &ToggleServiceUseCase{serviceRepo: repo}
}

func (uc *ToggleServiceUseCase) Execute(id uint) (*service.Service, error) {
	s, err := uc.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, service.ErrNotFound
	}
	s.Active = !s.Active
	if err := uc.serviceRepo.Update(s); err != nil {
		return nil, err
	}
	return s, nil
}
