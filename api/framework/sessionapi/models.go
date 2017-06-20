package sessionapi

// LoginModel is a binding model used for receiving the authentication data
type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
