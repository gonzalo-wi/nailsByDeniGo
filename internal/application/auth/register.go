package authapp

import "apiGoShei/internal/domain/client"

type RegisterInput struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

type RegisterOutput struct {
	ClientID uint
	Token    string
}

type RegisterUseCase struct {
	clientRepo client.Repository
	hasher     PasswordHasher
	tokenGen   TokenGenerator
}

func NewRegisterUseCase(repo client.Repository, hasher PasswordHasher, tokenGen TokenGenerator) *RegisterUseCase {
	return &RegisterUseCase{clientRepo: repo, hasher: hasher, tokenGen: tokenGen}
}

func (uc *RegisterUseCase) Execute(input RegisterInput) (*RegisterOutput, error) {
	existing, err := uc.clientRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, client.ErrEmailAlreadyExists
	}

	hash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	newClient := &client.Client{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: hash,
		Active:       true,
	}

	if err := uc.clientRepo.Create(newClient); err != nil {
		return nil, err
	}

	token, err := uc.tokenGen.Generate(newClient.ID, "client")
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{ClientID: newClient.ID, Token: token}, nil
}
