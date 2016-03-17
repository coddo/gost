package userloginservice

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"gost/service"
	"time"
)

const CollectionName = "user_sessions"

func CreateUserSession(userSession *dbmodels.UserSession) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if userSession.ID == "" {
		userSession.ID = bson.NewObjectId()
	}

	err := collection.Insert(userSession)

	return err
}

func UpdateUserSession(userSession *dbmodels.UserSession) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if userSession.ID == "" {
		return service.NoIdSpecifiedError
	}

	err := collection.UpdateId(userSession.ID, userSession)

	return err
}

func DeleteUserSession(sessionId bson.ObjectId) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	err := collection.RemoveId(sessionId)

	return err
}

func GetUserSession(token string) (*dbmodels.UserSession, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	userSession := dbmodels.UserSession{}
	err := collection.Find(bson.M{"token": token}).One(&userSession)

	return &userSession, err
}

func DeleteExpiredSessionsForUser(userId bson.ObjectId) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	_, err := collection.RemoveAll(bson.M{
		"userId":     userId,
		"expireDate": bson.M{"$lte": time.Now().Local()},
	})

	return err
}
