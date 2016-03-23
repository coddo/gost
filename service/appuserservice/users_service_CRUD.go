package appuserservice

import (
	"gost/orm/dbmodels"
	"gost/service"

	"gopkg.in/mgo.v2/bson"
)

const collectionName = "appusers"

// CreateUser adds a new ApplicationUser to the database
func CreateUser(user *dbmodels.ApplicationUser) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if user.ID == "" {
		user.ID = bson.NewObjectId()
	}

	err := collection.Insert(user)

	return err
}

// UpdateUser updates an existing ApplicationUser in the database
func UpdateUser(user *dbmodels.ApplicationUser) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if user.ID == "" {
		return service.ErrNoIDSpecified
	}

	err := collection.UpdateId(user.ID, user)

	return err
}

// GetUser retrieves an ApplicationUser from the database, based on its ID
func GetUser(userID bson.ObjectId) (*dbmodels.ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	user := dbmodels.ApplicationUser{}
	err := collection.FindId(userID).One(&user)

	return &user, err
}

// GetAllUsers retrieves all the existing ApplicationUser entities in the database
func GetAllUsers() ([]dbmodels.ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var users []dbmodels.ApplicationUser
	err := collection.Find(bson.M{}).All(&users)

	return users, err
}

// GetAllUsersLimited retrieves the first X ApplicationUser entities from the database, where X is the specified limit
func GetAllUsersLimited(limit int) ([]dbmodels.ApplicationUser, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var users []dbmodels.ApplicationUser
	err := collection.Find(bson.M{}).Limit(limit).All(&users)

	return users, err
}
