package authapi

import (
	"errors"
	"gost/api"
	"gost/util"
	"gost/util/jsonutil"
)

// Errors returned by the authapi
var (
	ErrInvalidEmail = errors.New("The email parameter was not specified or is invalid")
)

// RouteActivateAccount performs data parsing and binding before calling the API
func RouteActivateAccount(request *api.Request) api.Response {
	var model = &ActivateAccountModel{}

	var err = jsonutil.DeserializeJSON(request.Body, model)
	if err != nil {
		return api.BadRequest(err)
	}

	return activateAccount(model)
}

// RouteResendAccountActivationEmail performs data parsing and binding before calling the API
func RouteResendAccountActivationEmail(request *api.Request) api.Response {
	var email = request.GetStringRouteValue("email")
	if len(email) == 0 || !util.IsValidEmail(email) {
		return api.BadRequest(ErrInvalidEmail)
	}

	return resendAccountActivationEmail(email)
}

// RouteRequestResetPassword performs data parsing and binding before calling the API
func RouteRequestResetPassword(request *api.Request) api.Response {
	var email = request.GetStringRouteValue("email")
	if len(email) == 0 || !util.IsValidEmail(email) {
		return api.BadRequest(ErrInvalidEmail)
	}

	return requestResetPassword(email)
}

// RouteResetPassword performs data parsing and binding before calling the API
func RouteResetPassword(request *api.Request) api.Response {
	var model = &ResetPasswordModel{}

	var err = jsonutil.DeserializeJSON(request.Body, model)
	if err != nil {
		return api.BadRequest(err)
	}

	return resetPassword(model)
}

// RouteChangePassword performs data parsing and binding before calling the API
func RouteChangePassword(request *api.Request) api.Response {
	var model = &ChangePasswordModel{}

	var err = jsonutil.DeserializeJSON(request.Body, model)
	if err != nil {
		return api.BadRequest(err)
	}

	return changePassword(model)
}
