// Package service contains all the services necessary for interacting
// with all the collections(tables) in the database.
//
// Each service file should be written in it's own file and should represent only one dbmodels
package service

import (
	"errors"
	"gost/config"
	"log"

	"gopkg.in/mgo.v2"
)

var mongoDbSession *mgo.Session

var (
	// ErrNoIDSpecified says that no ID was specified for a fetch operation
	ErrNoIDSpecified = errors.New("No Id was specified for the entity")
)

// InitDbService initializes the connection (known as session) parameters to the database
func InitDbService() {
	if config.DbConnectionString == "" {
		log.Fatal("Database error: No connection string provided")
	}

	if mongoDbSession == nil {
		var err error
		mongoDbSession, err = mgo.Dial(config.DbConnectionString)
		mongoDbSession.SetMode(mgo.Monotonic, true)

		if err != nil {
			log.Fatalf("Can't connect to mongo, go error: %v\n", err)
		}
	}
}

// CloseDbService closes the current mongodb session
func CloseDbService() {
	mongoDbSession.Close()
}

// Connect creates the connection to the database.
// To avoid forgetting to close the session, always write: *defer session.Close()* right after getting it through this method
func Connect(collectionName string) (*mgo.Session, *mgo.Collection) {
	var session = mongoDbSession.Copy()

	return session, session.DB(config.DbName).C(collectionName)
}
