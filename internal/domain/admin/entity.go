package admin

import "time"

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "superadmin"
)

type Admin struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
