package models

import (
	"gost/dbmodels"
	"gost/service/appuserservice"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UserSession is a struct representing the session that a user has.
// Sessions are active since login until they expire or the user disconnects
type UserSession struct {
	ID bson.ObjectId `json:"id"`

	ApplicationUser ApplicationUser `json:"user"`
	Token           string          `json:"token"`
	ExpireDate      time.Time       `json:"expireDate"`
}

// PopConstrains fetches all the components from the database, based on their unique identifiers
func (userSession *UserSession) PopConstrains() {
	dbUser, err := appuserservice.GetUser(userSession.ApplicationUser.ID)
	if err == nil {
		userSession.ApplicationUser.Expand(dbUser)
	}
}

// Expand copies the dbmodels.UserSession to a UserSession expands all
// the components by fetching them from the database
func (userSession *UserSession) Expand(dbUserSession *dbmodels.UserSession) {
	userSession.ID = dbUserSession.ID
	userSession.ApplicationUser.ID = dbUserSession.UserID
	userSession.Token = dbUserSession.Token
	userSession.ExpireDate = dbUserSession.ExpireDate

	userSession.PopConstrains()
}

// Collapse coppies the UserSession to a dbmodels.UserSession user and
// only keeps the unique identifiers from the inner components
func (userSession *UserSession) Collapse() *dbmodels.UserSession {
	dbUserSession := dbmodels.UserSession{
		ID:         userSession.ID,
		UserID:     userSession.ApplicationUser.ID,
		Token:      userSession.Token,
		ExpireDate: userSession.ExpireDate,
	}

	return &dbUserSession
}
