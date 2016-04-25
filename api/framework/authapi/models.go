package authapi

// ActivateAccountModel is used for activating accounts
type ActivateAccountModel struct {
	Token string `json:"token"`
}

// RequestResetPasswordModel is used for requesting password resets over email
type RequestResetPasswordModel struct {
	Email                    string `json:"email"`
	PasswordResetServiceLink string `json:"passwordResetServiceLink"`
}

// ResetPasswordModel is used for resetting the account password
type ResetPasswordModel struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// ResendActivationEmailModel is used for resending the account activation email
type ResendActivationEmailModel struct {
	Email                      string `json:"email"`
	ActivateAccountServiceLink string `json:"activateAccountServiceLink"`
}
