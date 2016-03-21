package bll

import (
	"gost/api"
	"gost/dbmodels"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/appuserservice"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// GetApplicationUser retrieves a ApplicationUser
func GetApplicationUser(userID bson.ObjectId) api.Response {
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

// GetAllApplicationUsers retrieves all the existing ApplicationUser entities
func GetAllApplicationUsers() api.Response {
	var dbUsers []dbmodels.ApplicationUser

	dbUsers, err := appuserservice.GetAllUsers()
	return getAllApplicationUsers(dbUsers, err)
}

// GetAllApplicationUsersLimited retrieves a limited number of existing ApplicationUser entities
func GetAllApplicationUsersLimited(limit int) api.Response {
	var dbUsers []dbmodels.ApplicationUser

	dbUsers, err := appuserservice.GetAllUsersLimited(limit)
	return getAllApplicationUsers(dbUsers, err)
}

// CreateApplicationUser creates a new ApplicationUser
func CreateApplicationUser(user *models.ApplicationUser) api.Response {
	if !apifilter.CheckUserIntegrity(user) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	dbUser := user.Collapse()
	if dbUser == nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}

	err := appuserservice.CreateUser(dbUser)
	if err != nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}
	user.ID = dbUser.ID

	return api.SingleDataResponse(http.StatusCreated, user)
}

// UpdateApplicationUser updates an existing ApplicationUser
func UpdateApplicationUser(user *models.ApplicationUser) api.Response {
	if !apifilter.CheckUserIntegrity(user) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	dbUser := user.Collapse()
	if dbUser == nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}

	err := appuserservice.UpdateUser(dbUser)
	if err != nil {
		return api.NotFound(api.ErrEntityNotFound)
	}

	return api.SingleDataResponse(http.StatusOK, user)
}

func getAllApplicationUsers(dbUsers []dbmodels.ApplicationUser, err error) api.Response {
	if err != nil {
		return api.InternalServerError(err)
	}

	users := make([]*models.ApplicationUser, len(dbUsers))
	for i := 0; i < len(dbUsers); i++ {
		user := &models.ApplicationUser{}
		user.Expand(&dbUsers[i])

		users[i] = user
	}

	return api.MultipleDataResponse(http.StatusOK, users)
}
