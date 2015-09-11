package userservice

import (
	"gopkg.in/mgo.v2/bson"
	"gost/config"
	"gost/dbmodels"
	"gost/service"
	"testing"
)

func TestUserCRUD(t *testing.T) {
	user := &dbmodels.User{}

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
	config.InitTestsDatabase()
	service.InitDbService()

	if recover() != nil {
		t.Fatal("Test setup failed!")
	}
}

func tearDownUsersTest(t *testing.T, user *dbmodels.User) {
	err := DeleteUser(user.Id)

	if err != nil {
		t.Fatal("The user document could not be deleted!")
	}
}

func createUser(t *testing.T, user *dbmodels.User) {
	*user = dbmodels.User{
		Id:          bson.NewObjectId(),
		Password:    "CoddoPass",
		AccountType: dbmodels.AdministratorAccountType,
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

	err := CreateUser(user)

	if err != nil {
		t.Fatal("The user document could not be created!")
	}
}

func changeAndUpdateUser(t *testing.T, user *dbmodels.User) {
	user.Email = "testEmailCHanged@email.go"
	user.Password = "ChangedPassword"
	user.City = "Timisoara"
	user.PostalCode = 12512521
	user.AccountType = dbmodels.ClientAccountType

	err := UpdateUser(user)

	if err != nil {
		t.Fatal("The user document could not be updated!")
	}
}

func verifyUserCorresponds(t *testing.T, user *dbmodels.User) {
	dbuser, err := GetUser(user.Id)

	if err != nil || dbuser == nil {
		t.Error("Could not fetch the user document from the database!")
	}

	if !dbuser.Equal(user) {
		t.Error("The user document doesn't correspond with the document extracted from the database!")
	}
}
