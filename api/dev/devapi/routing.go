package devapi

import (
	"errors"
	"gost/api"
	"gost/util/jsonutil"
)

// Errors returned by the devapi
var (
	ErrInvalidToken   = errors.New("The token parameter is missing or invalid")
	ErrInvalidAppUser = errors.New("The application user data is missing or is invalid")
)

// RouteActivateAppUser performs data parsing and binding before calling the API
func RouteActivateAppUser(request *api.Request) api.Response {
	var token = request.GetStringParameter("token")
	if len(token) == 0 {
		return api.BadRequest(ErrInvalidToken)
	}

	return activateAppUser(token)
}

// RouteCreateAppUser performs data parsing and binding before calling the API
func RouteCreateAppUser(request *api.Request) api.Response {
	model := &AppUserModel{}

	err := jsonutil.DeserializeJSON(request.Body, model)
	if err != nil {
		return api.BadRequest(ErrInvalidAppUser)
	}

	return createAppUser(model)
}
