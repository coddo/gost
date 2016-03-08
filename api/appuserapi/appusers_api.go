package appuserapi

import (
	"gopkg.in/mgo.v2/bson"
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

func (usersApi *ApplicationUsersApi) GetUser(vars *api.ApiVar) api.ApiResponse {
	userId, err, found := apifilter.GetIdFromParams(vars.RequestForm)
	if found {
		if err != nil {
			return api.BadRequest(err)
		}

		return getUser(vars, userId)
	}

	limit, err, found := apifilter.GetIntValueFromParams("limit", vars.RequestForm)
	if found {
		if err != nil {
			return api.BadRequest(err)
		}

		return getAllUsers(vars, limit)
	}

	return getAllUsers(vars, -1)
}

func (usersApi *ApplicationUsersApi) CreateUser(vars *api.ApiVar) api.ApiResponse {
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

func (usersApi *ApplicationUsersApi) UpdateUser(vars *api.ApiVar) api.ApiResponse {
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

func getAllUsers(vars *api.ApiVar, limit int) api.ApiResponse {
	var dbUsers []dbmodels.ApplicationUser
	var err error

	if limit == 0 {
		return api.BadRequest(api.LimitParamError)
	}

	if limit != -1 {
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

func getUser(vars *api.ApiVar, userId bson.ObjectId) api.ApiResponse {
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
