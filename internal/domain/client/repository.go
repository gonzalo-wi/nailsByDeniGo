package client

type Repository interface {
	FindByID(id uint) (*Client, error)
	FindByEmail(email string) (*Client, error)
	FindAll() ([]Client, error)
	Create(client *Client) error
	Update(client *Client) error
}
