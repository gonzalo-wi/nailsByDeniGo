package models

import "gorm.io/gorm"

type AdminModel struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(20);default:'admin'"`
	Active       bool   `gorm:"default:true"`
}

func (AdminModel) TableName() string { return "admins" }
