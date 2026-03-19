package dto

type TokenResponse struct {
	Token     string `json:"token"`
	ClientID  uint   `json:"client_id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type AdminTokenResponse struct {
	Token   string `json:"token"`
	AdminID uint   `json:"admin_id"`
}
