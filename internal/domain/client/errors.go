package client

import "errors"

var (
	ErrNotFound           = errors.New("clienta no encontrada")
	ErrEmailAlreadyExists = errors.New("ya existe una cuenta con ese email")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
)
