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

type Service struct {
	session *mgo.Session
}

var mongoDBService *Service = &Service{}

var (
	NoIdSpecifiedError = errors.New("No Id was specified for the entity")
)

// Initialize a main session to the
func InitDbService() {
	url := config.DbConnectionString
	if url == "" {
		log.Fatal("Database error: No connection string provided")
	}

	if mongoDBService.session == nil {
		var err error
		mongoDBService.session, err = mgo.Dial(url)

		if err != nil {
			log.Fatalf("Can't connect to mongo, go error: %v\n", err)
		}
	}
}

// Close the mongodb service
func CloseDbService() {
	mongoDBService.session.Close()
}

// All connections should be stateless, so this method always
// returns a pointer to a new session. After each session is used,
// it should be closed in order to dump data from memory correctly.
// To avoid forgetting to close the session, always write: *defer session.Close()*
// right after getting it through this method
func Connect(collectionName string) (*mgo.Session, *mgo.Collection) {
	return mongoDBService.session.Copy(), mongoDBService.session.DB(config.DbName).C(collectionName)
}
