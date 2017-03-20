package authapi

import (
	"errors"
	"gost/api"
	"gost/auth"
	"gost/util/jsonutil"
	"net/http"
)

var errPasswordsDoNotMatch = errors.New("The password and its confirmation do not match")

// ActivateAccount activates an account using the activation token sent through email
func ActivateAccount(params *api.Request) api.Response {
	var model = ActivateAccountModel{}

	var err = jsonutil.DeserializeJSON(params.Body, &model)
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
func ResendAccountActivationEmail(params *api.Request) api.Response {
	var model = ResendActivationEmailModel{}

	var err = jsonutil.DeserializeJSON(params.Body, &model)
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
func RequestResetPassword(params *api.Request) api.Response {
	var model = RequestResetPasswordModel{}

	var err = jsonutil.DeserializeJSON(params.Body, &model)
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
func ResetPassword(params *api.Request) api.Response {
	var model = ResetPasswordModel{}

	var err = jsonutil.DeserializeJSON(params.Body, &model)
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

// ChangePassword changes an user account's password
func ChangePassword(params *api.Request) api.Response {
	var model = ChangePasswordModel{}

	var err = jsonutil.DeserializeJSON(params.Body, &model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if model.Password != model.PasswordConfirmation {
		return api.BadRequest(errPasswordsDoNotMatch)
	}

	err = auth.ChangePassword(model.Email, model.OldPassword, model.Password)
	if err != nil {
		return api.BadRequest(err)
	}

	return api.StatusResponse(http.StatusOK)
}
