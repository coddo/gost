package auth

import (
	"errors"
	"fmt"
	"gost/auth/identity"
	"gost/email"
	"gost/util"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	passwordResetTokenExpireTime     = 24 * time.Hour
	accountActivationTokenExpireTime = 7 * 24 * time.Hour
)

// Errors that can occur during ApplicationUser management
var (
	ErrActivationTokenExpired    = errors.New("The activation token has expired")
	ErrResetPasswordTokenExpired = errors.New("The reset password token has expired")
)

// CreateAppUser creates a new ApplicationUser with the given data, generates an activation token
// and sends an email containing a link used for activating the account
func CreateAppUser(emailAddress, password string, accountType int, activationServiceLink string) (*identity.ApplicationUser, error) {
	var token, err = util.GenerateUUID()
	if err != nil {
		return nil, err
	}

	passwordHash, err := util.HashString(password)
	if err != nil {
		return nil, err
	}

	var user = &identity.ApplicationUser{
		ID:                             bson.NewObjectId(),
		Email:                          emailAddress,
		Password:                       passwordHash,
		AccountType:                    accountType,
		ActivateAccountToken:           token,
		ActivateAccountTokenExpireDate: util.NextDateFromNow(accountActivationTokenExpireTime),
		AccountStatus:                  identity.AccountStatusDeactivated,
	}

	err = identity.CreateUser(user)
	if err != nil {
		return nil, err
	}

	go sendAccountActivationEmail(emailAddress, activationServiceLink, token)

	return user, nil
}

// ActivateAppUser activates an application user based on its token
func ActivateAppUser(token string) error {
	var user, err = identity.GetUserByActivationToken(token)
	if err != nil {
		return err
	}

	if util.IsDateExpiredFromNow(user.ActivateAccountTokenExpireDate) {
		return ErrActivationTokenExpired
	}

	user.AccountStatus = identity.AccountStatusActivated

	return identity.UpdateUser(user)
}

// ResetPassword resets the password of an application user
func ResetPassword(token, password string) error {
	var user, err = identity.GetUserByResetPasswordToken(token)
	if err != nil {
		return err
	}

	if util.IsDateExpiredFromNow(user.ResetPasswordTokenExpireDate) {
		return ErrResetPasswordTokenExpired
	}

	return changeUserPassword(user, password)
}

// ChangePassword changes the current password that the user has
func ChangePassword(userEmail, oldPassword, password string) error {
	var user, err = identity.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}

	if !util.MatchString(user.Password, oldPassword) {
		return ErrPasswordMismatch
	}

	return changeUserPassword(user, password)
}

// RequestResetPassword generates a reset token and sends an email with the link where to perform the change
func RequestResetPassword(emailAddress, passwordResetServiceLink string) error {
	var user, err = identity.GetUserByEmail(emailAddress)
	if err != nil {
		return err
	}

	token, err := util.GenerateUUID()
	if err != nil {
		return err
	}

	user.ResetPasswordToken = token
	user.ResetPasswordTokenExpireDate = util.NextDateFromNow(passwordResetTokenExpireTime)

	err = identity.UpdateUser(user)
	if err != nil {
		return err
	}

	go sendPasswordResetEmail(emailAddress, passwordResetServiceLink, token)

	return nil
}

// ResendAccountActivationEmail resends the email with the details for activating their user account
func ResendAccountActivationEmail(emailAddress, activationServiceLink string) error {
	var user, err = identity.GetUserByEmail(emailAddress)
	if err != nil {
		return err
	}

	token, err := util.GenerateUUID()
	if err != nil {
		return err
	}

	user.ActivateAccountToken = token
	user.ActivateAccountTokenExpireDate = util.NextDateFromNow(accountActivationTokenExpireTime)

	err = identity.UpdateUser(user)
	if err != nil {
		return err
	}

	go sendAccountActivationEmail(emailAddress, activationServiceLink, token)

	return nil
}

func sendAccountActivationEmail(userEmail, activationServiceLink, token string) {
	var accountActivationLink = fmt.Sprintf(activationServiceLink, token)

	err := email.SendAccountActivationEmail(userEmail, accountActivationLink)

	if err != nil {
		log.Printf(fmt.Sprintf("Error in sending account activation email to: %s", userEmail))
	}
}

func sendPasswordResetEmail(userEmail, passwordResetServiceLink, token string) {
	var passwordResetLink = fmt.Sprintf(passwordResetServiceLink, token)

	err := email.SendPasswordResetEmail(userEmail, passwordResetLink)

	if err != nil {
		log.Printf(fmt.Sprintf("Error in sending password reset email to: %s", userEmail))
	}
}

func changeUserPassword(user *identity.ApplicationUser, password string) error {
	passwordHash, err := util.HashString(password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	return identity.UpdateUser(user)
}
