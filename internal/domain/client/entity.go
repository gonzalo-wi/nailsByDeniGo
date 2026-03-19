package client

import "time"

type Client struct {
	ID           uint
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	PasswordHash string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
