package identity

import (
	"gost/dal/service"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Constants describing the status of the current user account
const (
	AccountStatusDeactivated = iota
	AccountStatusActivated   = iota
)

const collectionName = "appusers"

// ApplicationUser contains information necessary for managing accounts
type ApplicationUser struct {
	ID bson.ObjectId `bson:"_id" json:"id"`

	Email                          string    `bson:"email,omitempty" json:"email"`
	Password                       string    `bson:"password,omitempty" json:"password"`
	Roles                          []string  `bson:"roles,omitempty" json:"roles"`
	ResetPasswordToken             string    `bson:"resetPasswordToken,omitempty" json:"resetPasswordToken"`
	ResetPasswordTokenExpireDate   time.Time `bson:"resetPasswordTokenExpireDate,omitempty" json:"resetPasswordTokenExpireDate"`
	ActivateAccountToken           string    `bson:"activateAccountToken" json:"activateAccountToken"`
	ActivateAccountTokenExpireDate time.Time `bson:"activateAccountTokenExpireDate,omitempty" json:"activateAccountTokenExpireDate"`
	AccountStatus                  int       `bson:"accountStatus,omitempty" json:"accountStatus"`
}

// CreateUser adds a new ApplicationUser to the database
func CreateUser(user *ApplicationUser) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if user.ID == "" {
		user.ID = bson.NewObjectId()
	}

	var err = collection.Insert(user)

	return err
}

// UpdateUser updates an existing ApplicationUser in the database
func UpdateUser(user *ApplicationUser) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if user.ID == "" {
		return service.ErrNoIDSpecified
	}

	var err = collection.UpdateId(user.ID, user)

	return err
}

// GetUser retrieves an ApplicationUser from the database, based on its ID
func GetUser(userID bson.ObjectId) (*ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var user = ApplicationUser{}
	var err = collection.FindId(userID).One(&user)

	return &user, err
}

// GetUserByActivationToken retrieves an ApplicationUser from the database, based on its account activation token
func GetUserByActivationToken(token string) (*ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var user = ApplicationUser{}
	var err = collection.Find(bson.M{"activateAccountToken": token}).One(&user)

	return &user, err
}

// GetUserByResetPasswordToken retrieves an ApplicationUser from the database, based on its reset password token
func GetUserByResetPasswordToken(token string) (*ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var user = ApplicationUser{}
	var err = collection.Find(bson.M{"resetPasswordToken": token}).One(&user)

	return &user, err
}

// GetUserByEmail retrieves an ApplicationUser from the database, based on its email address
func GetUserByEmail(emailAddress string) (*ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var user = ApplicationUser{}
	var err = collection.Find(bson.M{"email": emailAddress}).One(&user)

	return &user, err
}

// DeleteUser deletes an ApplicationUser from the database, based on its ID
func DeleteUser(userID bson.ObjectId) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	return collection.RemoveId(userID)
}

// IsUserExistent verifies if an user with the given id exists
func IsUserExistent(userID bson.ObjectId) (*ApplicationUser, bool) {
	var user, err = GetUser(userID)

	return user, err == nil && user != nil
}

// IsUserEmailExistent verifies if an user with the given email exists
func IsUserEmailExistent(email string) (*ApplicationUser, bool) {
	var user, err = GetUserByEmail(email)

	return user, err == nil && user != nil
}

// IsUserActivated verifies if an user account is activated
func IsUserActivated(userID bson.ObjectId) (*ApplicationUser, bool) {
	user, exists := IsUserExistent(userID)
	if !exists || user == nil {
		return nil, false
	}

	return user, user.AccountStatus == AccountStatusActivated
}
