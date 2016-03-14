package appuserapi

import (
	"errors"
	"gost/api"
	"gost/dbmodels"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/appuserservice"
	"net/http"
	"strings"
)

// ApplicationUsersAPI defines the API endpoint for application user management
type ApplicationUsersAPI int

var (
	errLimitParam = errors.New("The limit cannot be 0. Use the value -1 for retrieving all the entities")
)

// Get endpoint fetches an application user based on a provided ID
func (usersApi *ApplicationUsersAPI) Get(vars *api.Request) api.Response {
	userID, err, found := apifilter.GetIdFromParams(vars.Form)

	if err != nil {
		return api.BadRequest(err)
	}

	if !found {
		return api.NotFound(err)
	}

	dbUser, err := appuserservice.GetUser(userID)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "found") {
			return api.NotFound(err)
		}

		return api.InternalServerError(err)
	}

	if dbUser == nil {
		return api.NotFound(api.ErrEntityNotFound)
	}

	user := &models.ApplicationUser{}
	user.Expand(dbUser)

	return api.SingleDataResponse(http.StatusOK, user)
}

// GetAll endpoint fetches all the existing users in the application.
// The number of returned entities can be limited using a "limit" request parameter
func (usersApi *ApplicationUsersAPI) GetAll(vars *api.Request) api.Response {
	limit, err, isLimitSpecified := apifilter.GetIntValueFromParams("limit", vars.Form)

	if err != nil {
		return api.BadRequest(err)
	}

	var dbUsers []dbmodels.ApplicationUser

	if isLimitSpecified {
		if limit == 0 {
			return api.BadRequest(errLimitParam)
		}

		dbUsers, err = appuserservice.GetAllUsersLimited(limit)
	} else {
		dbUsers, err = appuserservice.GetAllUsers()
	}

	if err != nil {
		return api.InternalServerError(err)
	}

	users := make([]models.Modeler, len(dbUsers))
	for i := 0; i < len(dbUsers); i++ {
		user := &models.ApplicationUser{}
		user.Expand(&dbUsers[i])

		users[i] = user
	}

	return api.MultipleDataResponse(http.StatusOK, users)
}

// Create endpoint creates a new application user
func (usersApi *ApplicationUsersAPI) Create(vars *api.Request) api.Response {
	user := &models.ApplicationUser{}

	err := models.DeserializeJson(vars.Body, user)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if !apifilter.CheckUserIntegrity(user) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	dbUser := user.Collapse()
	if dbUser == nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}

	err = appuserservice.CreateUser(dbUser)
	if err != nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}
	user.Id = dbUser.Id

	return api.SingleDataResponse(http.StatusCreated, user)
}

// Update endpoint updates an existing application user
func (usersApi *ApplicationUsersAPI) Update(vars *api.Request) api.Response {
	user := &models.ApplicationUser{}
	err := models.DeserializeJson(vars.Body, user)

	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if user.Id == "" {
		return api.BadRequest(api.ErrIDParamNotSpecified)
	}

	if !apifilter.CheckUserIntegrity(user) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	dbUser := user.Collapse()
	if dbUser == nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}

	err = appuserservice.UpdateUser(dbUser)
	if err != nil {
		return api.NotFound(api.ErrEntityNotFound)
	}

	return api.SingleDataResponse(http.StatusOK, user)
}
