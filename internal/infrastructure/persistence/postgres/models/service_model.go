package models

import "gorm.io/gorm"

type ServiceModel struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Description      string
	DurationMinutes  int     `gorm:"not null"`
	BasePrice        float64 `gorm:"type:numeric(10,2);not null"`
	RequiresDeposit  bool    `gorm:"default:false"`
	SuggestedDeposit float64 `gorm:"type:numeric(10,2);default:0"`
	Color            string  `gorm:"default:'#ffffff'"`
	Active           bool    `gorm:"default:true"`
}

func (ServiceModel) TableName() string { return "services" }
