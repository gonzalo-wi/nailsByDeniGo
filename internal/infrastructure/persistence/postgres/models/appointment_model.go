package models

import "gorm.io/gorm"

type AppointmentModel struct {
	gorm.Model
	ClientID       uint `gorm:"not null"`
	ServiceID      uint `gorm:"not null"`
	ProfessionalID *uint
	Date           string  `gorm:"type:varchar(10);not null"` // "2006-01-02"
	StartTime      string  `gorm:"type:varchar(5);not null"`  // "14:00"
	EndTime        string  `gorm:"type:varchar(5);not null"`
	BasePrice      float64 `gorm:"type:numeric(10,2);not null"`
	ExtrasAmount   float64 `gorm:"type:numeric(10,2);default:0"`
	ExtrasNote     string  `gorm:"type:text"`
	DepositAmount  float64 `gorm:"type:numeric(10,2);default:0"`
	FinalPrice     float64 `gorm:"type:numeric(10,2);not null"`
	Status        string  `gorm:"type:varchar(20);not null;default:'PENDING'"`
	Notes         string
	PenaltyAmount float64 `gorm:"type:numeric(10,2);default:0"`
	PenaltyNote   string  `gorm:"type:text"`

	Client  ClientModel  `gorm:"foreignKey:ClientID"`
	Service ServiceModel `gorm:"foreignKey:ServiceID"`
}

func (AppointmentModel) TableName() string { return "appointments" }
