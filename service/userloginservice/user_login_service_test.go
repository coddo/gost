package userloginservice

import (
	"gopkg.in/mgo.v2/bson"
	"gost/config"
	"gost/dbmodels"
	"gost/service"
	"testing"
	"time"
)

func TestUserSessionCRUD(t *testing.T) {
	userSession := &dbmodels.UserSession{}

	setUpUserSessionsTest(t)
	defer tearDownUserSessionsTest(t, userSession)

	testCreateUserSession(t, userSession)
	testVerifyUserSessionCorresponds(t, userSession)

	if t.Failed() {
		return
	}

	testChangeAndUpdateUserSession(t, userSession)
	testVerifyUserSessionCorresponds(t, userSession)

	testDeleteExpiredUserSessions(t)
}

func setUpUserSessionsTest(t *testing.T) {
	config.InitTestsDatabase()
	service.InitDbService()

	if recover() != nil {
		t.Fatal("Test setup failed!")
	}
}

func tearDownUserSessionsTest(t *testing.T, userSession *dbmodels.UserSession) {
	err := DeleteUserSession(userSession.Id)

	if err != nil {
		t.Fatal("The user session document could not be deleted!")
	}
}

func testCreateUserSession(t *testing.T, userSession *dbmodels.UserSession) {
	*userSession = dbmodels.UserSession{
		UserId:     bson.NewObjectId(),
		Token:      "trh46rth46rth4r",
		ExpireDate: time.Now().Local(),
	}

	err := CreateUserSession(userSession)

	if err != nil {
		t.Fatal("The user session document could not be created!")
	}
}

func testChangeAndUpdateUserSession(t *testing.T, userSession *dbmodels.UserSession) {
	userSession.UserId = bson.NewObjectId()
	userSession.Token = "a65g4as65as4g6as4ga"
	userSession.ExpireDate = time.Date(2015, time.December, 12, 0, 0, 0, 0, time.UTC)

	err := UpdateUserSession(userSession)

	if err != nil {
		t.Fatal("The user session document could not be updated!")
	}
}

func testVerifyUserSessionCorresponds(t *testing.T, userSession *dbmodels.UserSession) {
	dbUserSession, err := GetUserSession(userSession.Token)

	if err != nil || dbUserSession == nil {
		t.Error("Could not fetch the user session document from the database!")
	}

	if !dbUserSession.Equal(userSession) {
		t.Error("The user session document doesn't correspond with the document extracted from the database!")
	}
}

func testDeleteExpiredUserSessions(t *testing.T) {
	userSession1 := &dbmodels.UserSession{
		Id:         bson.NewObjectId(),
		UserId:     bson.NewObjectId(),
		Token:      "as7f6as8faf5aasf6721rqf",
		ExpireDate: time.Now().Local().Add(-time.Hour * 150),
	}

	userSession2 := &dbmodels.UserSession{
		Id:         bson.NewObjectId(),
		UserId:     userSession1.UserId,
		Token:      "a68f4asg6546sgafas4f6a",
		ExpireDate: time.Now().Local().Add(-time.Hour * 300),
	}

	err1 := CreateUserSession(userSession1)
	err2 := CreateUserSession(userSession2)
	if err1 != nil || err2 != nil {
		t.Fatal("Error creating expired user sessions!")
	}

	DeleteExpiredSessionsForUser(userSession1.UserId)
	_, err1 = GetUserSession(userSession1.Token)
	_, err2 = GetUserSession(userSession2.Token)

	if err1 == nil && err2 == nil {
		t.Fatal("The expired user sessions haven't been properly deleted!")
	}
}
