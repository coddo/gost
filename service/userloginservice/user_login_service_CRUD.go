package userloginservice

import (
	"gost/dbmodels"
	"gost/service"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const collectionName = "user_sessions"

// CreateUserSession adds a new UserSession to the database
func CreateUserSession(userSession *dbmodels.UserSession) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if userSession.ID == "" {
		userSession.ID = bson.NewObjectId()
	}

	err := collection.Insert(userSession)

	return err
}

// UpdateUserSession updates an existing UserSession in the database
func UpdateUserSession(userSession *dbmodels.UserSession) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if userSession.ID == "" {
		return service.ErrNoIDSpecified
	}

	err := collection.UpdateId(userSession.ID, userSession)

	return err
}

// DeleteUserSession removes an UserSession from the database
func DeleteUserSession(sessionID bson.ObjectId) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	err := collection.RemoveId(sessionID)

	return err
}

// GetUserSession retrieves an UserSession from the database, based on its ID
func GetUserSession(token string) (*dbmodels.UserSession, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	userSession := dbmodels.UserSession{}
	err := collection.Find(bson.M{"token": token}).One(&userSession)

	return &userSession, err
}

// DeleteExpiredSessionsForUser removes all the user sessions that have expired from the database
func DeleteExpiredSessionsForUser(userID bson.ObjectId) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	_, err := collection.RemoveAll(bson.M{
		"userId":     userID,
		"expireDate": bson.M{"$lte": time.Now().Local()},
	})

	return err
}
