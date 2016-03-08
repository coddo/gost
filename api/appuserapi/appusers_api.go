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

type ApplicationUsersApi int

const ApiName = "applicationUsers"

var (
	limitParamError = errors.New("The limit cannot be 0. Use the value -1 for retrieving all the entities")
)

func (usersApi *ApplicationUsersApi) Get(vars *api.ApiVar) api.ApiResponse {
	userId, err, found := apifilter.GetIdFromParams(vars.RequestForm)

	if err != nil {
		return api.BadRequest(err)
	}

	if !found {
		return api.NotFound(err)
	}

	dbUser, err := appuserservice.GetUser(userId)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "found") {
			return api.NotFound(err)
		} else {
			return api.InternalServerError(err)
		}
	}

	if dbUser == nil {
		return api.NotFound(api.EntityNotFoundError)
	}

	user := &models.ApplicationUser{}
	user.Expand(dbUser)

	return api.SingleDataResponse(http.StatusOK, user)
}

func (usersApi *ApplicationUsersApi) GetAll(vars *api.ApiVar) api.ApiResponse {
	limit, err, isLimitSpecified := apifilter.GetIntValueFromParams("limit", vars.RequestForm)

	if err != nil {
		return api.BadRequest(err)
	}

	var dbUsers []dbmodels.ApplicationUser

	if isLimitSpecified {
		if limit == 0 {
			return api.BadRequest(limitParamError)
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

func (usersApi *ApplicationUsersApi) Create(vars *api.ApiVar) api.ApiResponse {
	user := &models.ApplicationUser{}

	err := models.DeserializeJson(vars.RequestBody, user)
	if err != nil {
		return api.BadRequest(api.EntityFormatError)
	}

	if !apifilter.CheckUserIntegrity(user) {
		return api.BadRequest(api.EntityIntegrityError)
	}

	dbUser := user.Collapse()
	if dbUser == nil {
		return api.InternalServerError(api.EntityProcessError)
	}

	err = appuserservice.CreateUser(dbUser)
	if err != nil {
		return api.InternalServerError(api.EntityProcessError)
	}
	user.Id = dbUser.Id

	return api.SingleDataResponse(http.StatusCreated, user)
}

func (usersApi *ApplicationUsersApi) Update(vars *api.ApiVar) api.ApiResponse {
	user := &models.ApplicationUser{}
	err := models.DeserializeJson(vars.RequestBody, user)

	if err != nil {
		return api.BadRequest(api.EntityFormatError)
	}

	if user.Id == "" {
		return api.BadRequest(api.IdParamNotSpecifiedError)
	}

	if !apifilter.CheckUserIntegrity(user) {
		return api.BadRequest(api.EntityIntegrityError)
	}

	dbUser := user.Collapse()
	if dbUser == nil {
		return api.InternalServerError(api.EntityProcessError)
	}

	err = appuserservice.UpdateUser(dbUser)
	if err != nil {
		return api.NotFound(api.EntityNotFoundError)
	}

	return api.SingleDataResponse(http.StatusOK, user)
}
