package repositories

import (
	"errors"

	"apiGoShei/internal/domain/service"
	"apiGoShei/internal/infrastructure/persistence/postgres/models"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) FindByID(id uint) (*service.Service, error) {
	var m models.ServiceModel
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toServiceDomain(&m), nil
}

func (r *ServiceRepository) FindAll(activeOnly bool) ([]service.Service, error) {
	var rows []models.ServiceModel
	q := r.db
	if activeOnly {
		q = q.Where("active = ?", true)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]service.Service, len(rows))
	for i, m := range rows {
		result[i] = *toServiceDomain(&m)
	}
	return result, nil
}

func (r *ServiceRepository) Create(s *service.Service) error {
	m := &models.ServiceModel{
		Name:             s.Name,
		Description:      s.Description,
		DurationMinutes:  s.DurationMinutes,
		BasePrice:        s.BasePrice,
		RequiresDeposit:  s.RequiresDeposit,
		SuggestedDeposit: s.SuggestedDeposit,
		Color:            s.Color,
		Active:           s.Active,
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	s.ID = m.ID
	s.CreatedAt = m.CreatedAt
	s.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *ServiceRepository) Update(s *service.Service) error {
	return r.db.Model(&models.ServiceModel{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
		"name":              s.Name,
		"description":       s.Description,
		"duration_minutes":  s.DurationMinutes,
		"base_price":        s.BasePrice,
		"requires_deposit":  s.RequiresDeposit,
		"suggested_deposit": s.SuggestedDeposit,
		"color":             s.Color,
		"active":            s.Active,
	}).Error
}

func toServiceDomain(m *models.ServiceModel) *service.Service {
	return &service.Service{
		ID:               m.ID,
		Name:             m.Name,
		Description:      m.Description,
		DurationMinutes:  m.DurationMinutes,
		BasePrice:        m.BasePrice,
		RequiresDeposit:  m.RequiresDeposit,
		SuggestedDeposit: m.SuggestedDeposit,
		Color:            m.Color,
		Active:           m.Active,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
