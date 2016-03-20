package userloginapi

import (
	"fmt"
	"gost/api"
	"gost/models"
	"gost/service/userloginservice"
	"gost/tests"
	"net/http"
	"net/url"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	GET    = "Get"
	CREATE = "Create"
)

const apiPath = "/users/login"

var userSessionsRoute = fmt.Sprintf(`[{"id": "UserSessionsRoute", "pattern": "/users/login", 
    "handlers": {"%s": "GET", "%s": "POST"}}]`, GET, CREATE)

type dummyUserSession struct {
	BadField string
}

func (userSession *dummyUserSession) PopConstrains() {}

func TestUserSessionsApi(t *testing.T) {
	tests.InitializeServerConfigurations(userSessionsRoute, new(UserSessionsAPI))

	testCreateUserSessionInBadFormat(t)
	sessionID, token := testCreateUserSessionInGoodFormat(t)
	testGetUserSessionWithInexistentTokenInDB(t)
	testGetUserSessionWithGoodIDParam(t, token)

	userloginservice.DeleteUserSession(sessionID)
}

func testGetUserSessionWithInexistentTokenInDB(t *testing.T) {
	params := url.Values{}
	params.Add("token", "asagasgsaga7615651")

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetUserSessionWithGoodIDParam(t *testing.T, token string) {
	params := url.Values{}
	params.Add("token", token)

	rw := tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

}

func testCreateUserSessionInBadFormat(t *testing.T) {
	dUserSession := &dummyUserSession{
		BadField: "bad value",
	}

	tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusBadRequest, nil, dUserSession, t)
}

func testCreateUserSessionInGoodFormat(t *testing.T) (bson.ObjectId, string) {
	userSession := &models.UserSession{
		ID:              bson.NewObjectId(),
		ApplicationUser: models.ApplicationUser{ID: bson.NewObjectId()},
		Token:           "as7f6as8faf5aasf6721rqf",
		ExpireDate:      time.Now().Local(),
	}

	rw := tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusCreated, nil, userSession, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return userSession.ID, userSession.Token
}
