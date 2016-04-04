package valuesapi

import (
	"gost/api"
	"gost/auth/identity"
	"gost/util"
	"net/http"
)

// CreateAppUser is an endpoint used for creating application users
func (v *ValuesAPI) CreateAppUser(params *api.Request) api.Response {
	user := &identity.ApplicationUser{}

	err := util.DeserializeJSON(params.Body, user)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	err = identity.CreateUser(user)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.StatusResponse(http.StatusOK)
}
