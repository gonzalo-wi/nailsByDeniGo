package authapp

import (
	"errors"

	"apiGoShei/internal/domain/admin"
)

var ErrAdminInvalidCredentials = errors.New("credenciales inválidas")

type AdminLoginInput struct {
	Email    string
	Password string
}

type AdminLoginOutput struct {
	AdminID uint
	Token   string
}

type AdminLoginUseCase struct {
	adminRepo admin.Repository
	hasher    PasswordHasher
	tokenGen  TokenGenerator
}

func NewAdminLoginUseCase(repo admin.Repository, hasher PasswordHasher, tokenGen TokenGenerator) *AdminLoginUseCase {
	return &AdminLoginUseCase{adminRepo: repo, hasher: hasher, tokenGen: tokenGen}
}

func (uc *AdminLoginUseCase) Execute(input AdminLoginInput) (*AdminLoginOutput, error) {
	a, err := uc.adminRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if a == nil || !a.Active || !uc.hasher.Check(input.Password, a.PasswordHash) {
		return nil, ErrAdminInvalidCredentials
	}

	token, err := uc.tokenGen.Generate(a.ID, string(a.Role))
	if err != nil {
		return nil, err
	}

	return &AdminLoginOutput{AdminID: a.ID, Token: token}, nil
}
