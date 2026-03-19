package authapp

// PasswordHasher abstrae la operación de hashing de contraseñas.
// La implementación concreta vive en infrastructure/security.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}

// TokenGenerator abstrae la generación de tokens de autenticación.
// La implementación concreta vive en infrastructure/security.
type TokenGenerator interface {
	Generate(userID uint, role string) (string, error)
}
