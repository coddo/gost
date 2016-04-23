package email

// SendAccountActivationEmail sends an account activation confirmation email to a user
func SendAccountActivationEmail(email, urlEndpoint string) error {
	emailBody := ParseTemplate(activateAccountTemplate, urlEndpoint)

	mail := NewEmail()
	mail.SetRecipient(email)
	mail.SetSubject(activateAccountSubject)
	mail.SetBody(emailBody)

	return mail.Send()
}

// SendPasswordResetEmail sends an account password reset instructions email to a user
func SendPasswordResetEmail(email, urlEndpoint string) error {
	emailBody := ParseTemplate(resetPasswordTemplate, urlEndpoint)

	mail := NewEmail()
	mail.SetRecipient(email)
	mail.SetSubject(resetPasswordSubject)
	mail.SetBody(emailBody)

	return mail.Send()
}
