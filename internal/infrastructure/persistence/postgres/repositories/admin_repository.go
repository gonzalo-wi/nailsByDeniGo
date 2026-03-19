package repositories

import (
	"errors"

	"apiGoShei/internal/domain/admin"
	"apiGoShei/internal/infrastructure/persistence/postgres/models"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByID(id uint) (*admin.Admin, error) {
	var m models.AdminModel
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toAdminDomain(&m), nil
}

func (r *AdminRepository) FindByEmail(email string) (*admin.Admin, error) {
	var m models.AdminModel
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toAdminDomain(&m), nil
}

func (r *AdminRepository) Create(a *admin.Admin) error {
	m := &models.AdminModel{
		Name:         a.Name,
		Email:        a.Email,
		PasswordHash: a.PasswordHash,
		Role:         string(a.Role),
		Active:       a.Active,
	}
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	a.ID = m.ID
	a.CreatedAt = m.CreatedAt
	a.UpdatedAt = m.UpdatedAt
	return nil
}

func toAdminDomain(m *models.AdminModel) *admin.Admin {
	return &admin.Admin{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         admin.Role(m.Role),
		Active:       m.Active,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
