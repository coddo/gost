package userapi

import (
	"gopkg.in/mgo.v2/bson"
	"gost/api"
	"gost/dbmodels"
	"gost/models"
	"gost/tests"
	"net/http"
	"net/url"
	"testing"
)

const usersRoute = "[{\"id\": \"UsersRoute\", \"pattern\": \"/users\", \"handlers\": {\"DELETE\": \"DeleteUser\", \"GET\": \"GetUser\", \"POST\": \"PostUser\", \"PUT\": \"PutUser\"}}]"
const apiPath = "/users"

type dummyUser struct {
	BadField string
}

func (user *dummyUser) PopConstrains() {}

func TestUsersApi(t *testing.T) {
	tests.InitializeServerConfigurations(usersRoute, new(UsersApi))

	testPostUserInBadFormat(t)
	id := testPostUserInGoodFormat(t)
	testPutUserInBadFormat(t)
	testPutUserWithoutId(t)
	testPutUserWithNoExistentIdInDb(t)
	testPutUserWithGoodRequestDetails(t, id)
	testGetUserWithInexistentIdInDB(t)
	testGetUserWithBadIdParam(t)
	testGetUserWithGoodIdParam(t, id)
	testGetAllUsersWithoutLimit(t)
	testGetAllUsersWithBadLimitParam(t)
	testGetAllUsersWithGoodLimitParam(t)
	testDeleteUserWithNoIdParam(t)
	testDeleteUserWithIdParamInWrongFormat(t)
	testDeteleUserWithInexistentIdInDB(t)
	testDeteleUserWithGoodRequestParams(t, id)
}

func testGetUserWithInexistentIdInDB(t *testing.T) {
	params := url.Values{}
	params.Add("id", bson.NewObjectId().Hex())

	tests.PerformApiTestCall(apiPath, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetUserWithBadIdParam(t *testing.T) {
	params := url.Values{}
	params.Add("id", "2as456fas4")

	tests.PerformApiTestCall(apiPath, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetUserWithGoodIdParam(t *testing.T, id bson.ObjectId) {
	params := url.Values{}
	params.Add("id", id.Hex())

	rw := tests.PerformApiTestCall(apiPath, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

}

func testGetAllUsersWithoutLimit(t *testing.T) {
	rw := tests.PerformApiTestCall(apiPath, api.GET, http.StatusOK, nil, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}
}

func testGetAllUsersWithBadLimitParam(t *testing.T) {
	params := url.Values{}
	params.Add("limit", "asfsa")

	tests.PerformApiTestCall(apiPath, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetAllUsersWithGoodLimitParam(t *testing.T) {
	params := url.Values{}
	params.Add("limit", "20")

	rw := tests.PerformApiTestCall(apiPath, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}
}

func testPostUserInBadFormat(t *testing.T) {
	dUser := &dummyUser{
		BadField: "bad value",
	}

	tests.PerformApiTestCall(apiPath, api.POST, http.StatusBadRequest, nil, dUser, t)
}

func testPostUserInGoodFormat(t *testing.T) bson.ObjectId {
	user := &models.User{
		Id:          bson.NewObjectId(),
		Password:    "CoddoPass",
		AccountType: dbmodels.ADMINISTRATOR_ACCOUNT_TYPE,
		FirstName:   "Claudiu",
		LastName:    "Codoban",
		Email:       "test@tests.com",
		Sex:         'M',
		Country:     "Romania",
		State:       "Hunedoara",
		City:        "Deva",
		Address:     "AddrTest",
		PostalCode:  330099,
		Picture:     "ftp://pictLink",
		Token:       "as7f6as8faf5aasf6721rqf",
	}

	rw := tests.PerformApiTestCall(apiPath, api.POST, http.StatusCreated, nil, user, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return user.Id
}

func testPutUserInBadFormat(t *testing.T) {
	user := &models.User{
		Id:        "507f191e810c19729de860ea",
		Sex:       'M',
		FirstName: "gigel",
		Country:   "Romania",
	}

	tests.PerformApiTestCall(apiPath, api.PUT, http.StatusBadRequest, nil, user, t)
}

func testPutUserWithoutId(t *testing.T) {
	user := &models.User{
		Email:     "ceva@ceva.com",
		Sex:       'M',
		FirstName: "gigel",
		Token:     "fsa4fas564g6g4s6ag",
		Country:   "Romania",
	}

	tests.PerformApiTestCall(apiPath, api.PUT, http.StatusBadRequest, nil, user, t)
}

func testPutUserWithNoExistentIdInDb(t *testing.T) {
	user := &models.User{
		Id:        bson.NewObjectId(),
		Email:     "ceva@ceva.com",
		Sex:       'M',
		FirstName: "gigel",
		Token:     "fsa4fas564g6g4s6ag",
		Country:   "Romania",
		Address:   "addr",
	}

	tests.PerformApiTestCall(apiPath, api.PUT, http.StatusNotFound, nil, user, t)
}

func testPutUserWithGoodRequestDetails(t *testing.T, id bson.ObjectId) {
	user := &models.User{
		Id:        id,
		Email:     "ceva@ceva.com",
		Sex:       'M',
		FirstName: "gigel",
		Token:     "fsa4fas564g6g4s6ag",
		Country:   "Romania",
		Address:   "addr",
	}

	rw := tests.PerformApiTestCall(apiPath, api.PUT, http.StatusOK, nil, user, t)
	body := rw.Body.String()

	if len(body) == 0 {
		t.Fatal("The response body was wither empty or deteriorated", body)
	}
}

func testDeleteUserWithNoIdParam(t *testing.T) {
	tests.PerformApiTestCall(apiPath, api.DELETE, http.StatusBadRequest, nil, nil, t)
}

func testDeleteUserWithIdParamInWrongFormat(t *testing.T) {
	params := url.Values{}
	params.Add("id", "a46fsa65gas")

	tests.PerformApiTestCall(apiPath, api.DELETE, http.StatusBadRequest, params, nil, t)
}

func testDeteleUserWithInexistentIdInDB(t *testing.T) {
	params := url.Values{}
	params.Add("id", bson.NewObjectId().Hex())

	tests.PerformApiTestCall(apiPath, api.DELETE, http.StatusNotFound, params, nil, t)
}

func testDeteleUserWithGoodRequestParams(t *testing.T, id bson.ObjectId) {
	params := url.Values{}
	params.Add("id", id.Hex())

	tests.PerformApiTestCall(apiPath, api.DELETE, http.StatusNoContent, params, nil, t)
}
