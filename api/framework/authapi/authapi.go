package authapi

import (
	"errors"
	"gost/api"
	"gost/auth"
	"gost/util"
	"net/http"
)

// AuthAPI defines the API endpoint for user authorization
type AuthAPI int

var errPasswordsDoNotMatch = errors.New("The password and its confirmation do not match")

// ActivateAccount activates an account using the activation token sent through email
func (a *AuthAPI) ActivateAccount(params *api.Request) api.Response {
	var model = ActivateAccountModel{}

	var err = util.DeserializeJSON(params.Body, &model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	err = auth.ActivateAppUser(model.Token)
	if err != nil {
		return api.BadRequest(err)
	}

	return api.StatusResponse(http.StatusOK)
}

// ResendAccountActivationEmail resends the email with the details for activating their user account
func (a *AuthAPI) ResendAccountActivationEmail(params *api.Request) api.Response {
	var model = ResendActivationEmailModel{}

	var err = util.DeserializeJSON(params.Body, &model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	err = auth.ResendAccountActivationEmail(model.Email, model.ActivateAccountServiceLink)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.StatusResponse(http.StatusOK)
}

// RequestResetPassword sends an email with a special token that will be used for resetting the password
func (a *AuthAPI) RequestResetPassword(params *api.Request) api.Response {
	var model = RequestResetPasswordModel{}

	var err = util.DeserializeJSON(params.Body, &model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	err = auth.RequestResetPassword(model.Email, model.PasswordResetServiceLink)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.StatusResponse(http.StatusOK)
}

// ResetPassword resets an user account's password
func (a *AuthAPI) ResetPassword(params *api.Request) api.Response {
	var model = ResetPasswordModel{}

	var err = util.DeserializeJSON(params.Body, &model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if model.Password != model.PasswordConfirmation {
		return api.BadRequest(errPasswordsDoNotMatch)
	}

	err = auth.ResetPassword(model.Token, model.Password)
	if err != nil {
		if err == auth.ErrResetPasswordTokenExpired {
			return api.BadRequest(err)
		}

		return api.InternalServerError(err)
	}

	return api.StatusResponse(http.StatusOK)
}
