package appuserservice

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"gost/service"
)

const CollectionName = "users"

func CreateUser(user *dbmodels.ApplicationUser) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if user.ID == "" {
		user.ID = bson.NewObjectId()
	}

	err := collection.Insert(user)

	return err
}

func UpdateUser(user *dbmodels.ApplicationUser) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if user.ID == "" {
		return service.NoIdSpecifiedError
	}

	err := collection.UpdateId(user.ID, user)

	return err
}

func GetUser(userId bson.ObjectId) (*dbmodels.ApplicationUser, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	user := dbmodels.ApplicationUser{}
	err := collection.FindId(userId).One(&user)

	return &user, err
}

func GetAllUsers() ([]dbmodels.ApplicationUser, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	var users []dbmodels.ApplicationUser
	err := collection.Find(bson.M{}).All(&users)

	return users, err
}

func GetAllUsersLimited(limit int) ([]dbmodels.ApplicationUser, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	var users []dbmodels.ApplicationUser
	err := collection.Find(bson.M{}).Limit(limit).All(&users)

	return users, err
}
