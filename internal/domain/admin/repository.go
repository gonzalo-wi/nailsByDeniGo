package admin

type Repository interface {
	FindByID(id uint) (*Admin, error)
	FindByEmail(email string) (*Admin, error)
	Create(admin *Admin) error
}
