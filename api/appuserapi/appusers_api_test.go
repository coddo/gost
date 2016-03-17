package appuserapi

import (
	"fmt"
	"gost/api"
	"gost/dbmodels"
	"gost/models"
	"gost/tests"
	"net/http"
	"net/url"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

const (
	GET    = "Get"
	GETALL = "GetAll"
	CREATE = "Create"
	UPDATE = "Update"
)

const apiPath = "/appusers"

var applicationUsersRoute = fmt.Sprintf(`[{"id": "ApplicationUsersRoute", "pattern": "/appusers", 
    "handlers": {"%s": "GET", "%s": "GET", "%s": "POST", "%s": "PUT"}}]`, GET, GETALL, CREATE, UPDATE)

type dummyUser struct {
	BadField string
}

func (user *dummyUser) PopConstrains() {}

func TestUsersApi(t *testing.T) {
	tests.InitializeServerConfigurations(applicationUsersRoute, new(ApplicationUsersAPI))

	testCreateUserInBadFormat(t)
	id := testCreateUserInGoodFormat(t)
	testUpdateUserInBadFormat(t)
	testUpdateUserWithoutID(t)
	testUpdateUserWithNoExistentIDInDb(t)
	testUpdateUserWithGoodRequestDetails(t, id)
	testGetUserWithInexistentIDInDatabase(t)
	testGetUserWithBadIDParam(t)
	testGetUserWithGoodIDParam(t, id)
	testGetAllUsersWithoutLimit(t)
	testGetAllUsersWithBadLimitParam(t)
	testGetAllUsersWithZeroLimitParam(t)
	testGetAllUsersWithGoodLimitParam(t)
}

func testGetUserWithInexistentIDInDatabase(t *testing.T) {
	params := url.Values{}
	params.Add("id", bson.NewObjectId().Hex())

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetUserWithBadIDParam(t *testing.T) {
	params := url.Values{}
	params.Add("id", "2as456fas4")

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetUserWithGoodIDParam(t *testing.T, id bson.ObjectId) {
	params := url.Values{}
	params.Add("id", id.Hex())

	rw := tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}
}

func testGetAllUsersWithoutLimit(t *testing.T) {
	rw := tests.PerformApiTestCall(apiPath, GETALL, api.GET, http.StatusOK, nil, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}
}

func testGetAllUsersWithBadLimitParam(t *testing.T) {
	params := url.Values{}
	params.Add("limit", "asfsa")

	tests.PerformApiTestCall(apiPath, GETALL, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetAllUsersWithZeroLimitParam(t *testing.T) {
	params := url.Values{}
	params.Add("limit", "0")

	tests.PerformApiTestCall(apiPath, GETALL, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetAllUsersWithGoodLimitParam(t *testing.T) {
	params := url.Values{}
	params.Add("limit", "20")

	rw := tests.PerformApiTestCall(apiPath, GETALL, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}
}

func testCreateUserInBadFormat(t *testing.T) {
	dUser := &dummyUser{
		BadField: "bad value",
	}

	tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusBadRequest, nil, dUser, t)
}

func testCreateUserInGoodFormat(t *testing.T) bson.ObjectId {
	user := &models.ApplicationUser{
		Id:                 bson.NewObjectId(),
		Password:           "CoddoPass",
		AccountType:        dbmodels.AdministratorAccountType,
		Email:              "test@tests.com",
		ResetPasswordToken: "as7f6as8faf5aasf6721rqf",
	}

	rw := tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusCreated, nil, user, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return user.Id
}

func testUpdateUserInBadFormat(t *testing.T) {
	user := &models.ApplicationUser{
		Id:                 "507f191e810c19729de860ea",
		ResetPasswordToken: "asg1a89wqg4a5s",
	}

	tests.PerformApiTestCall(apiPath, UPDATE, api.PUT, http.StatusBadRequest, nil, user, t)
}

func testUpdateUserWithoutID(t *testing.T) {
	user := &models.ApplicationUser{
		Email:              "ceva@ceva.com",
		Password:           "CoddoPass",
		ResetPasswordToken: "fsa4fas564g6g4s6ag",
	}

	tests.PerformApiTestCall(apiPath, UPDATE, api.PUT, http.StatusBadRequest, nil, user, t)
}

func testUpdateUserWithNoExistentIDInDb(t *testing.T) {
	user := &models.ApplicationUser{
		Id:                 bson.NewObjectId(),
		Email:              "ceva@ceva.com",
		Password:           "CoddoPass",
		ResetPasswordToken: "fsa4fas564g6g4s6ag",
	}

	tests.PerformApiTestCall(apiPath, UPDATE, api.PUT, http.StatusNotFound, nil, user, t)
}

func testUpdateUserWithGoodRequestDetails(t *testing.T, id bson.ObjectId) {
	user := &models.ApplicationUser{
		Id:                 id,
		Email:              "ceva@ceva.com",
		Password:           "CoddoPass",
		ResetPasswordToken: "fsa4fas564g6g4s6ag",
	}

	rw := tests.PerformApiTestCall(apiPath, UPDATE, api.PUT, http.StatusOK, nil, user, t)
	body := rw.Body.String()

	if len(body) == 0 {
		t.Fatal("The response body was wither empty or deteriorated", body)
	}
}
