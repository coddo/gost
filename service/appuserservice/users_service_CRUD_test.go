package appuserservice

import (
	"gost/dbmodels"
	"gost/service"
	testconfig "gost/tests/config"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestUserCRUD(t *testing.T) {
	user := &dbmodels.ApplicationUser{}

	setUpUsersTest(t)
	defer tearDownUsersTest(t, user)

	createUser(t, user)
	verifyUserCorresponds(t, user)

	if !t.Failed() {
		changeAndUpdateUser(t, user)
		verifyUserCorresponds(t, user)
	}
}

func setUpUsersTest(t *testing.T) {
	testconfig.InitTestsDatabase()
	service.InitDbService()

	if recover() != nil {
		t.Fatal("Test setup failed!")
	}
}

func tearDownUsersTest(t *testing.T, user *dbmodels.ApplicationUser) {
	err := deleteUser(user.Id)

	if err != nil {
		t.Fatal("The user document could not be deleted!")
	}
}

func createUser(t *testing.T, user *dbmodels.ApplicationUser) {
	*user = dbmodels.ApplicationUser{
		Id:                           bson.NewObjectId(),
		Password:                     "CoddoPass",
		AccountType:                  dbmodels.ADMINISTRATOR_ACCOUNT_TYPE,
		Email:                        "test@tests.com",
		ResetPasswordToken:           "as7f6as8faf5aasf6721rqf",
		ResetPasswordTokenExpireDate: time.Now(),
		Status: dbmodels.STATUS_ACCOUNT_ACTIVATED,
	}

	err := CreateUser(user)

	if err != nil {
		t.Fatal("The user document could not be created!")
	}
}

func changeAndUpdateUser(t *testing.T, user *dbmodels.ApplicationUser) {
	user.Email = "testEmailCHanged@email.go"
	user.Password = "ChangedPassword"
	user.AccountType = dbmodels.NORMAL_USER_ACCOUNT_TYPE
	user.Status = dbmodels.STATUS_ACCOUNT_DEACTIVATED

	err := UpdateUser(user)

	if err != nil {
		t.Fatal("The user document could not be updated!")
	}
}

func verifyUserCorresponds(t *testing.T, user *dbmodels.ApplicationUser) {
	dbuser, err := GetUser(user.Id)

	if err != nil || dbuser == nil {
		t.Error("Could not fetch the user document from the database!")
	}

	if !dbuser.Equal(user) {
		t.Error("The user document doesn't correspond with the document extracted from the database!")
	}
}

func deleteUser(userId bson.ObjectId) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	err := collection.RemoveId(userId)

	return err
}
