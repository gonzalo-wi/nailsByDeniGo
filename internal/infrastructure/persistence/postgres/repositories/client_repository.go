package repositories

import (
	"errors"

	"apiGoShei/internal/domain/client"
	"apiGoShei/internal/infrastructure/persistence/postgres/models"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) FindByID(id uint) (*client.Client, error) {
	var m models.ClientModel
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toClientDomain(&m), nil
}

func (r *ClientRepository) FindByEmail(email string) (*client.Client, error) {
	var m models.ClientModel
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toClientDomain(&m), nil
}

func (r *ClientRepository) FindAll() ([]client.Client, error) {
	var models []models.ClientModel
	if err := r.db.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]client.Client, len(models))
	for i, m := range models {
		m := m
		result[i] = *toClientDomain(&m)
	}
	return result, nil
}

func (r *ClientRepository) Create(c *client.Client) error {
	m := &models.ClientModel{
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Email:        c.Email,
		Phone:        c.Phone,
		PasswordHash: c.PasswordHash,
		Active:       c.Active,
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	c.ID = m.ID
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *ClientRepository) Update(c *client.Client) error {
	return r.db.Model(&models.ClientModel{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"email":      c.Email,
		"phone":      c.Phone,
		"active":     c.Active,
	}).Error
}

func toClientDomain(m *models.ClientModel) *client.Client {
	return &client.Client{
		ID:           m.ID,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		Email:        m.Email,
		Phone:        m.Phone,
		PasswordHash: m.PasswordHash,
		Active:       m.Active,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
