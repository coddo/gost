package identity

import (
	"gost/orm/service"
	testconfig "gost/tests/config"
	"gost/util/dateutil"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestUserCRUD(t *testing.T) {
	user := &ApplicationUser{}

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

func tearDownUsersTest(t *testing.T, user *ApplicationUser) {
	err := deleteUser(user.ID)

	if err != nil {
		t.Fatal("The user document could not be deleted!")
	}
}

func createUser(t *testing.T, user *ApplicationUser) {
	*user = ApplicationUser{
		ID:                           bson.NewObjectId(),
		Password:                     "CoddoPass",
		AccountType:                  AccountTypeAdministrator,
		Email:                        "test@tests.com",
		ResetPasswordToken:           "as7f6as8faf5aasf6721rqf",
		ResetPasswordTokenExpireDate: time.Now(),
		AccountStatus:                AccountStatusActivated,
	}

	err := CreateUser(user)

	if err != nil {
		t.Fatal("The user document could not be created!")
	}
}

func changeAndUpdateUser(t *testing.T, user *ApplicationUser) {
	user.Email = "testEmailCHanged@email.go"
	user.Password = "ChangedPassword"
	user.AccountType = AccountTypeNormalUser
	user.AccountStatus = AccountStatusDeactivated

	err := UpdateUser(user)

	if err != nil {
		t.Fatal("The user document could not be updated!")
	}
}

func verifyUserCorresponds(t *testing.T, user *ApplicationUser) {
	dbuser, err := GetUser(user.ID)

	if err != nil || dbuser == nil {
		t.Error("Could not fetch the user document from the database!")
	}

	if user.AccountStatus != dbuser.AccountStatus || user.AccountType != dbuser.AccountType ||
		user.ActivateAccountToken != dbuser.ActivateAccountToken ||
		!dateutil.CompareDates(user.ActivateAccountTokenExpireDate, dbuser.ActivateAccountTokenExpireDate) ||
		user.Email != dbuser.Email || user.Password != dbuser.Password ||
		user.ResetPasswordToken != dbuser.ResetPasswordToken ||
		!dateutil.CompareDates(user.ResetPasswordTokenExpireDate, dbuser.ResetPasswordTokenExpireDate) {

		t.Error("The user document doesn't correspond with the document extracted from the database!")
	}
}

func deleteUser(userID bson.ObjectId) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	err := collection.RemoveId(userID)

	return err
}
