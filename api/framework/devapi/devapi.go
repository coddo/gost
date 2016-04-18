package devapi

import (
	"gost/api"
	"gost/auth/identity"
	"gost/util"
	"net/http"
)

// DevAPI defines the API endpoint for development actions and custom testing
type DevAPI int

// CreateAppUser is an endpoint used for creating application users
func (v *DevAPI) CreateAppUser(params *api.Request) api.Response {
	user := &identity.ApplicationUser{}

	err := util.DeserializeJSON(params.Body, user)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	err = identity.CreateUser(user)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.JSONResponse(http.StatusOK, user)
}
