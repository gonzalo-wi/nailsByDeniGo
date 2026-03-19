package service

type Repository interface {
	FindByID(id uint) (*Service, error)
	FindAll(activeOnly bool) ([]Service, error)
	Create(s *Service) error
	Update(s *Service) error
}
