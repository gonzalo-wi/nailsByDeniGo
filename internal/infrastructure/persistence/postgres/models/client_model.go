package models

import "gorm.io/gorm"

type ClientModel struct {
	gorm.Model
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	Phone        string
	PasswordHash string `gorm:"not null"`
	Active       bool   `gorm:"default:true"`
}

func (ClientModel) TableName() string { return "clients" }
