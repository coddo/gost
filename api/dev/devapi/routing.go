package devapi

import (
	"gost/api"
	"gost/util/jsonutil"
)

// RouteActivateAppUser performs data parsing and binding before calling the API
func RouteActivateAppUser(request *api.Request) api.Response {
	var token, found = api.GetStringParameter("token", request.Form)
	if !found {
		return api.BadRequest(api.ErrInvalidInput)
	}

	return activateAppUser(token)
}

// RouteCreateAppUser performs data parsing and binding before calling the API
func RouteCreateAppUser(request *api.Request) api.Response {
	model := &AppUserModel{}
	err := jsonutil.DeserializeJSON(request.Body, model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return createAppUser(model)
}
