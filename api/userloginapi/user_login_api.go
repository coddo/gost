package userloginapi

import (
	"errors"
	"gost/api"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/userloginservice"
	"gost/util"
	"net/http"
	"time"
)

// UserSessionsAPI defines the API endpoint for user sessions
type UserSessionsAPI int

var (
	// ErrTokenNotSpecified tells the user the token parameter was not specified
	ErrTokenNotSpecified = errors.New("The session token hasn't been specified")

	// ErrTokenExpired tells the user that the given token has expired
	ErrTokenExpired = errors.New("The session with the specified token has expired")
)

const (
	daysUntilExpire = 7
)

// Get endpoint retrieves a user session entity based on a given token
func (userSessionsApi *UserSessionsAPI) Get(vars *api.Request) api.Response {
	token, found := apifilter.GetStringValueFromParams("token", vars.Form)

	if !found {
		return api.BadRequest(ErrTokenNotSpecified)
	}

	dbUserSession, err := userloginservice.GetUserSession(token)
	if err != nil {
		return api.NotFound(err)
	} else if util.IsDateExpiredFromNow(dbUserSession.ExpireDate) {
		userloginservice.DeleteUserSession(dbUserSession.ID)
		return api.Unauthorized(ErrTokenExpired)
	}
	dbUserSession.ExpireDate = util.NextDateFromNow(time.Hour * 24 * daysUntilExpire)

	err = userloginservice.UpdateUserSession(dbUserSession)
	if err != nil {
		return api.InternalServerError(err)
	}

	userSession := new(models.UserSession)
	userSession.Expand(dbUserSession)

	userloginservice.DeleteExpiredSessionsForUser(dbUserSession.UserID)
	return api.SingleDataResponse(http.StatusOK, userSession)
}

// Create endpoint creates a new user session
func (userSessionsApi *UserSessionsAPI) Create(vars *api.Request) api.Response {
	userSession := &models.UserSession{}

	err := models.DeserializeJSON(vars.Body, userSession)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if !apifilter.CheckUserSessionIntegrity(userSession) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	userSession.ExpireDate = util.NextDateFromNow(time.Hour * 24 * daysUntilExpire)

	dbUserSession := userSession.Collapse()
	if dbUserSession == nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}

	err = userloginservice.CreateUserSession(dbUserSession)
	if err != nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}
	userSession.Id = dbUserSession.ID

	userloginservice.DeleteExpiredSessionsForUser(dbUserSession.UserID)
	return api.SingleDataResponse(http.StatusCreated, userSession)
}
