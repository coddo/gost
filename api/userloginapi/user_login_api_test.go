package userloginapi

import (
	"go-server-template/api"
	"go-server-template/models"
	"go-server-template/service/userloginservice"
	"go-server-template/tests"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const userSessionsRoute = "[{\"id\": \"UserSessionsRoute\", \"pattern\": \"/users/login\", \"handlers\": {\"DELETE\": \"DeleteUserSession\", \"GET\": \"GetUserSession\", \"POST\": \"PostUserSession\", \"PUT\": \"PutUserSession\"}}]"
const apiPath = "/users/login"

type dummyUserSession struct {
	BadField string
}

func (userSession *dummyUserSession) PopConstrains() {}

func TestUserSessionsApi(t *testing.T) {
	tests.InitializeServerConfigurations(userSessionsRoute, new(UserSessionsApi))

	testPostUserSessionInBadFormat(t)
	sessionId, token := testPostUserSessionInGoodFormat(t)
	testGetUserSessionWithInexistentTokenInDB(t)
	testGetUserSessionWithGoodIdParam(t, token)

	userloginservice.DeleteUserSession(sessionId)
}

func testGetUserSessionWithInexistentTokenInDB(t *testing.T) {
	params := url.Values{}
	params.Add("token", "asagasgsaga7615651")

	tests.PerformApiTestCall(apiPath, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetUserSessionWithGoodIdParam(t *testing.T, token string) {
	params := url.Values{}
	params.Add("token", token)

	rw := tests.PerformApiTestCall(apiPath, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

}

func testPostUserSessionInBadFormat(t *testing.T) {
	dUserSession := &dummyUserSession{
		BadField: "bad value",
	}

	tests.PerformApiTestCall(apiPath, api.POST, http.StatusBadRequest, nil, dUserSession, t)
}

func testPostUserSessionInGoodFormat(t *testing.T) (bson.ObjectId, string) {
	userSession := &models.UserSession{
		Id:         bson.NewObjectId(),
		User:       models.User{Id: bson.NewObjectId()},
		Token:      "as7f6as8faf5aasf6721rqf",
		ExpireDate: time.Now().Local(),
	}

	rw := tests.PerformApiTestCall(apiPath, api.POST, http.StatusCreated, nil, userSession, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return userSession.Id, userSession.Token
}
