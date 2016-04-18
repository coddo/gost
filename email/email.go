package email

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"
)

const (
	authServer   = "smtp.zoho.com"
	smtpServer   = "smtp.zoho.com:587"
	gostUsername = "gostwebframework@zoho.com"
	gostPassword = "gostwebframework"
	senderName   = "GostWebFramework"
	senderEmail  = "gostwebframework@zoho.com"
)

var (
	sender        = mail.Address{Name: senderName, Address: senderEmail}
	authorization = smtp.PlainAuth("", gostUsername, gostPassword, authServer)
)

var basicHeader = map[string]string{
	"From":                      sender.String(),
	"MIME-Version":              "1.0",
	"Content-Type":              "text/html; charset=\"utf-8\"",
	"Content-Transfer-Encoding": "base64",
}

// Email struct is used to send an email message
type Email struct {
	recipient []string
	body      string
	header    map[string]string
}

// NewEmail creates a new email message
func NewEmail() *Email {
	return &Email{
		header: basicHeader,
	}
}

// Send sends the email message
func (email *Email) Send() error {
	return smtp.SendMail(smtpServer,
		authorization,
		sender.Address,
		email.recipient,
		email.createContent())
}

// SetRecipient sets the receiver of the email
func (email *Email) SetRecipient(name, address string) {
	recipient := mail.Address{
		Name:    name,
		Address: address,
	}

	email.recipient = []string{recipient.Address}
	email.header["To"] = recipient.String()
}

// SetSubject sets the subject of the email
func (email *Email) SetSubject(subject string) {
	email.header["Subject"] = subject
}

// SetBody sets the body of the email
func (email *Email) SetBody(body string) {
	email.body = body
}

func (email *Email) createContent() []byte {
	var message bytes.Buffer

	// Header
	for key, value := range email.header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Body delimiter
	message.WriteString("\r\n")

	// Body
	message.WriteString(email.body)

	return message.Bytes()
}
