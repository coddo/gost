package authapi

import "gost/api"

// AuthAPI defines the API endpoint for user authorization
type AuthAPI int

// ActivateAccount activates an account using the activation token sent through email
func (a *AuthAPI) ActivateAccount(params *api.Request) api.Response {
	return api.Response{}
}

// RequestResetPassword sends an email with a special token that will be used for resetting the password
func (a *AuthAPI) RequestResetPassword(params *api.Request) api.Response {
	return api.Response{}
}

// ResetPassword resets an user account's password
func (a *AuthAPI) ResetPassword(params *api.Request) api.Response {
	if !params.Identity.IsAuthorized() {
		return api.Unauthorized()
	}

	return api.Response{}
}
