package authapp

import "apiGoShei/internal/domain/client"

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	ClientID  uint
	Token     string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type LoginUseCase struct {
	clientRepo client.Repository
	hasher     PasswordHasher
	tokenGen   TokenGenerator
}

func NewLoginUseCase(repo client.Repository, hasher PasswordHasher, tokenGen TokenGenerator) *LoginUseCase {
	return &LoginUseCase{clientRepo: repo, hasher: hasher, tokenGen: tokenGen}
}

func (uc *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	c, err := uc.clientRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if c == nil || !c.Active || !uc.hasher.Check(input.Password, c.PasswordHash) {
		return nil, client.ErrInvalidCredentials
	}

	token, err := uc.tokenGen.Generate(c.ID, "client")
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		ClientID:  c.ID,
		Token:     token,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
	}, nil
}
