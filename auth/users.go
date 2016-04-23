package auth

import (
	"errors"
	"fmt"
	"gost/auth/identity"
	"gost/email"
	"gost/util"
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
	}

	err = identity.CreateUser(user)
	if err != nil {
		return nil, err
	}

	var emailActivationLink = fmt.Sprintf(activationServiceLink, token)
	err = email.SendAccountActivationEmail(user.Email, emailActivationLink)
	if err != nil {
		rollbackAppUserCreation(user.ID)
		return nil, err
	}

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

	passwordHash, err := util.HashString(password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	return identity.UpdateUser(user)
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

	var passwordResetLink = fmt.Sprintf(passwordResetServiceLink, token)

	return email.SendPasswordResetEmail(user.Email, passwordResetLink)
}

func rollbackAppUserCreation(userID bson.ObjectId) {
	identity.DeleteUser(userID)
}
