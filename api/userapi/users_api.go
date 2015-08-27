package userapi

import (
	"gopkg.in/mgo.v2/bson"
	"gost/api"
	"gost/dbmodels"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/userservice"
	"net/http"
	"strings"
)

type UsersApi int

const ApiName = "users"

func (usersApi *UsersApi) GetUser(vars *api.ApiVar) api.ApiResponse {
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

func (usersApi *UsersApi) PostUser(vars *api.ApiVar) api.ApiResponse {
	user := &models.User{}

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

	err = userservice.CreateUser(dbUser)
	if err != nil {
		return api.InternalServerError(api.EntityProcessError)
	}
	user.Id = dbUser.Id

	return api.SingleDataResponse(http.StatusCreated, user)
}

func (usersApi *UsersApi) PutUser(vars *api.ApiVar) api.ApiResponse {
	user := &models.User{}
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

	err = userservice.UpdateUser(dbUser)
	if err != nil {
		return api.NotFound(api.EntityNotFoundError)
	}

	return api.SingleDataResponse(http.StatusOK, user)
}

func (usersApi *UsersApi) DeleteUser(vars *api.ApiVar) api.ApiResponse {
	userId, err, found := apifilter.GetIdFromParams(vars.RequestForm)

	if found {
		if err != nil {
			return api.BadRequest(err)
		}

		err = userservice.DeleteUser(userId)
		if err != nil {
			return api.NotFound(err)
		}

		return api.StatusResponse(http.StatusOK)
	}

	return api.BadRequest(err)
}

func getAllUsers(vars *api.ApiVar, limit int) api.ApiResponse {
	var dbUsers []dbmodels.User
	var err error

	if limit == 0 {
		return api.BadRequest(api.LimitParamError)
	}

	if limit != -1 {
		dbUsers, err = userservice.GetAllUsersLimited(limit)
	} else {
		dbUsers, err = userservice.GetAllUsers()
	}

	if err != nil {
		return api.InternalServerError(err)
	}

	usersMap := make(map[int]models.User, len(dbUsers))
	for i := 0; i < len(dbUsers); i++ {
		user := &models.User{}

		user.Expand(&dbUsers[i])

		usersMap[i] = *user
	}

	return api.MultipleDataResponse(http.StatusOK, usersMap)
}

func getUser(vars *api.ApiVar, userId bson.ObjectId) api.ApiResponse {
	dbUser, err := userservice.GetUser(userId)

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

	user := &models.User{}
	user.Expand(dbUser)

	return api.SingleDataResponse(http.StatusOK, user)
}
