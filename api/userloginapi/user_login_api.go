package userloginapi

import (
	"errors"
	"gost/api"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/userloginservice"
	"net/http"
	"time"
)

type UserSessionsApi int

const ApiName = "userSessions"

var (
	TokenNotSpecifiedError = errors.New("The session token hasn't been specified")
	TokenExpiredError      = errors.New("The session with the specified token has expired")
)

const (
	daysUntilExpire = 7
)

func (userSessionsApi *UserSessionsApi) GetUserSession(vars *api.ApiVar) api.ApiResponse {
	token, found := apifilter.GetStringValueFromParams("token", vars.RequestForm)
	today := time.Now().Local()

	if !found {
		return api.BadRequest(TokenNotSpecifiedError)
	}

	dbUserSession, err := userloginservice.GetUserSession(token)
	if err != nil {
		return api.NotFound(err)
	} else if dbUserSession.ExpireDate.Local().Before(today) {
		userloginservice.DeleteUserSession(dbUserSession.Id)
		return api.Unauthorized(TokenExpiredError)
	}
	dbUserSession.ExpireDate = today.Add(time.Hour * 24 * daysUntilExpire)

	err = userloginservice.UpdateUserSession(dbUserSession)
	if err != nil {
		return api.InternalServerError(err)
	}

	userSession := new(models.UserSession)
	userSession.Expand(dbUserSession)

	userloginservice.DeleteExpiredSessionsForUser(dbUserSession.UserId)
	return api.SingleDataResponse(http.StatusOK, userSession)
}

func (userSessionsApi *UserSessionsApi) PostUserSession(vars *api.ApiVar) api.ApiResponse {
	userSession := &models.UserSession{}

	err := models.DeserializeJson(vars.RequestBody, userSession)
	if err != nil {
		return api.BadRequest(api.EntityFormatError)
	}

	if !apifilter.CheckUserSessionIntegrity(userSession) {
		return api.BadRequest(api.EntityIntegrityError)
	}

	today := time.Now().Local()
	userSession.ExpireDate = today.Add(time.Hour * 24 * daysUntilExpire)

	dbUserSession := userSession.Collapse()
	if dbUserSession == nil {
		return api.InternalServerError(api.EntityProcessError)
	}

	err = userloginservice.CreateUserSession(dbUserSession)
	if err != nil {
		return api.InternalServerError(api.EntityProcessError)
	}
	userSession.Id = dbUserSession.Id

	userloginservice.DeleteExpiredSessionsForUser(dbUserSession.UserId)
	return api.SingleDataResponse(http.StatusCreated, userSession)
}
