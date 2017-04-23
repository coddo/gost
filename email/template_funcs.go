package email

// SendAccountActivationEmail sends an account activation confirmation email to a user
func SendAccountActivationEmail(emailAddress, urlEndpoint string) error {
	var email = NewEmail(emailAddress, activateAccountSubject, activateAccountTemplate, urlEndpoint)

	return email.Send()
}

// SendPasswordResetEmail sends an account password reset instructions email to a user
func SendPasswordResetEmail(emailAddress, urlEndpoint string) error {
	var email = NewEmail(emailAddress, resetPasswordSubject, resetPasswordTemplate, urlEndpoint)

	return email.Send()
}
