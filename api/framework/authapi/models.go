package authapi

// ActivateAccountModel is used for activating accounts
type ActivateAccountModel struct {
	Token string `json:"token"`
}

// RequestResetPassword is used for requesting password resets over email
type RequestResetPassword struct {
	Email                    string `json:"email"`
	PasswordResetServiceLink string `json:"passwordResetServiceLink"`
}

// ResetPassword is used for resetting the account password
type ResetPassword struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}
