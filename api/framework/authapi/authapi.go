package authapi

import (
	"errors"
	"gost/api"
	"gost/auth"
	"net/http"
)

// Errors returned by the authapi
var (
	ErrPasswordsDoNotMatch = errors.New("The password and its confirmation do not match")
)

// activateAccount activates an account using the activation token sent through email
func activateAccount(model *ActivateAccountModel) api.Response {
	var err = auth.ActivateAppUser(model.Token)

	if err != nil {
		return api.BadRequest(err)
	}

	return api.StatusResponse(http.StatusNoContent)
}

// resendAccountActivationEmail resends the email with the details for activating their user account
func resendAccountActivationEmail(email string) api.Response {
	activationLink, err := auth.ResendAccountActivationEmail(email)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.PlainTextResponse(http.StatusOK, activationLink)
}

// requestResetPassword sends an email with a special token that will be used for resetting the password
func requestResetPassword(email string) api.Response {
	var resetLink, err = auth.RequestResetPassword(email)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.PlainTextResponse(http.StatusOK, resetLink)
}

// resetPassword resets an user account's password
func resetPassword(model *ResetPasswordModel) api.Response {
	if model.Password != model.PasswordConfirmation {
		return api.BadRequest(ErrPasswordsDoNotMatch)
	}

	var err = auth.ResetPassword(model.Token, model.Password)
	if err != nil {
		if err == auth.ErrResetPasswordTokenExpired {
			return api.BadRequest(err)
		}

		return api.InternalServerError(err)
	}

	return api.StatusResponse(http.StatusNoContent)
}

// changePassword changes an user account's password
func changePassword(model *ChangePasswordModel) api.Response {
	if model.Password != model.PasswordConfirmation {
		return api.BadRequest(ErrPasswordsDoNotMatch)
	}

	var err = auth.ChangePassword(model.Email, model.OldPassword, model.Password)
	if err != nil {
		return api.BadRequest(err)
	}

	return api.StatusResponse(http.StatusNoContent)
}
