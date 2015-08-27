//  Package which contains all the services necessary for interacting
//  with all the collections(tables) in the database.
//
// Each service file should be written in it's own file and should represent only one dbmodels
package service

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gost/config"
	"log"
)

var (
	NoIdSpecifiedError = errors.New("No Id was specified for the entity")
)

// All connections should be stateless, so this method always
// returns a pointer to a new session. After each session is used,
// it should be closed in order to dump data from memory correctly.
// To avoid forgetting to close the session, always write: *defer session.Close()*
// right after getting it through this method
func Connect(collectionName string) (*mgo.Session, *mgo.Collection) {
	uri := config.DbConnectionString
	if uri == "" {
		log.Fatal("Database error: No connection string provided")
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		log.Fatalf("Can't connect to mongo, go error: %v\n", err)
	}

	sess.SetSafe(&mgo.Safe{})

	return sess, sess.DB(config.DbName).C(collectionName)
}
