package identity

import (
	"gost/service"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Constants describing the type of account the the users have
const (
	AccountTypeNormalUser    = iota
	AccountTypeAdministrator = iota
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
	AccountType                    int       `bson:"accountType,omitempty" json:"accountType"`
	ResetPasswordToken             string    `bson:"resetPasswordToken,omitempty" json:"resetPasswordToken"`
	ResetPasswordTokenExpireDate   time.Time `bson:"resetPasswordTokenExpireDate,omitempty" json:"resetPasswordTokenExpireDate"`
	ActivateAccountToken           string    `bson:"activateAccountToken" json:"activateAccountToken"`
	ActivateAccountTokenExpireDate time.Time `bson:"activateAccountTokenExpireDate,omitempty" json:"activateAccountTokenExpireDate"`
	AccountStatus                  int       `bson:"status,omitempty" json:"status"`
}

// CreateUser adds a new ApplicationUser to the database
func CreateUser(user *ApplicationUser) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if user.ID == "" {
		user.ID = bson.NewObjectId()
	}

	err := collection.Insert(user)

	return err
}

// UpdateUser updates an existing ApplicationUser in the database
func UpdateUser(user *ApplicationUser) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if user.ID == "" {
		return service.ErrNoIDSpecified
	}

	err := collection.UpdateId(user.ID, user)

	return err
}

// GetUser retrieves an ApplicationUser from the database, based on its ID
func GetUser(userID bson.ObjectId) (*ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	user := ApplicationUser{}
	err := collection.FindId(userID).One(&user)

	return &user, err
}

// Exists verifies if an user with the given id exists
func IsUserExistent(userID bson.ObjectId) bool {
	user, err := GetUser(userID)

	return err == nil && user != nil
}
