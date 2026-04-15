package authapp

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}

type TokenGenerator interface {
	Generate(userID uint, role string) (string, error)
}
