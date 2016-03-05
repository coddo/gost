package userservice

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"gost/service"
)

const CollectionName = "users"

func CreateUser(user *dbmodels.User) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if user.Id == "" {
		user.Id = bson.NewObjectId()
	}

	err := collection.Insert(user)

	return err
}

func UpdateUser(user *dbmodels.User) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if user.Id == "" {
		return service.NoIdSpecifiedError
	}

	err := collection.UpdateId(user.Id, user)

	return err
}

func DeleteUser(userId bson.ObjectId) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	err := collection.RemoveId(userId)

	return err
}

func GetUser(userId bson.ObjectId) (*dbmodels.User, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	user := dbmodels.User{}
	err := collection.FindId(userId).One(&user)

	return &user, err
}

func GetAllUsers() ([]dbmodels.User, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	var users []dbmodels.User
	err := collection.Find(bson.M{}).All(&users)

	return users, err
}

func GetAllUsersLimited(limit int) ([]dbmodels.User, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	var users []dbmodels.User
	err := collection.Find(bson.M{}).Limit(limit).All(&users)

	return users, err
}
