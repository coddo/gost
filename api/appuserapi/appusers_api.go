package appuserapi

import (
	"errors"
	"gost/api"
	"gost/bll"
	"gost/filter/apifilter"
	"gost/models"
	"gost/util"
)

// ApplicationUsersAPI defines the API endpoint for application user management
type ApplicationUsersAPI int

var (
	errLimitParam = errors.New("The limit cannot be 0. Use the value -1 for retrieving all the entities")
)

// Get endpoint fetches an application user based on a provided ID
func (usersApi *ApplicationUsersAPI) Get(vars *api.Request) api.Response {
	userID, found, err := apifilter.GetIDFromParams(vars.Form)

	if err != nil {
		return api.BadRequest(err)
	}

	if !found {
		api.BadRequest(api.ErrIDParamNotSpecified)
	}

	return bll.GetApplicationUser(userID)
}

// GetAll endpoint fetches all the existing users in the application.
// The number of returned entities can be limited using a "limit" request parameter
func (usersApi *ApplicationUsersAPI) GetAll(vars *api.Request) api.Response {
	limit, isLimitSpecified, err := apifilter.GetIntValueFromParams("limit", vars.Form)

	if err != nil {
		return api.BadRequest(err)
	}

	if !isLimitSpecified {
		return bll.GetAllApplicationUsers()
	}

	if limit == 0 {
		return api.BadRequest(errLimitParam)
	}

	return bll.GetAllApplicationUsersLimited(limit)
}

// Create endpoint creates a new application user
func (usersApi *ApplicationUsersAPI) Create(vars *api.Request) api.Response {
	user := &models.ApplicationUser{}

	err := util.DeserializeJSON(vars.Body, user)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return bll.CreateApplicationUser(user)
}

// Update endpoint updates an existing application user
func (usersApi *ApplicationUsersAPI) Update(vars *api.Request) api.Response {
	user := &models.ApplicationUser{}

	err := util.DeserializeJSON(vars.Body, user)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if user.ID == "" {
		return api.BadRequest(api.ErrIDParamNotSpecified)
	}

	return bll.UpdateApplicationUser(user)
}
