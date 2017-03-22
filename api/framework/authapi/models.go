package authapi

// ActivateAccountModel is used for activating accounts
type ActivateAccountModel struct {
	Token string `json:"token"`
}

// ResetPasswordModel is used for resetting the account password
type ResetPasswordModel struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// ChangePasswordModel is used for changing a user's password
type ChangePasswordModel struct {
	Email                string `json:"email"`
	OldPassword          string `json:"oldPassword"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}
