package serviceapp

import "apiGoShei/internal/domain/service"

type ListServiceUseCase struct {
	serviceRepo service.Repository
}

func NewListServiceUseCase(repo service.Repository) *ListServiceUseCase {
	return &ListServiceUseCase{serviceRepo: repo}
}

func (uc *ListServiceUseCase) Execute(activeOnly bool) ([]service.Service, error) {
	return uc.serviceRepo.FindAll(activeOnly)
}

// ─── GetByID ──────────────────────────────────────────────────────────────────

type GetServiceByIDUseCase struct {
	serviceRepo service.Repository
}

func NewGetServiceByIDUseCase(repo service.Repository) *GetServiceByIDUseCase {
	return &GetServiceByIDUseCase{serviceRepo: repo}
}

func (uc *GetServiceByIDUseCase) Execute(id uint) (*service.Service, error) {
	s, err := uc.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, service.ErrNotFound
	}
	return s, nil
}
