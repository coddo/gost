package dbmodels

import (
	"gost/util"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UserSession is a struct representing the session that a user has.
// Sessions are active since login until they expire or the user disconnects
type UserSession struct {
	ID bson.ObjectId `bson:"_id" json:"id"`

	UserID     bson.ObjectId `bson:"userId,omitempty" json:"userId"`
	Token      string        `bson:"token,omitempty" json:"token"`
	ExpireDate time.Time     `bson:"expireDate,omitempty" json:"expireDate"`
}

// Equal compares two UserSession objects. Implements the Objecter interface
func (userSession *UserSession) Equal(obj Objecter) bool {
	otherSession, ok := obj.(*UserSession)
	if !ok {
		return false
	}

	switch {
	case userSession.ID != otherSession.ID:
		return false
	case userSession.Token != otherSession.Token:
		return false
	case userSession.UserID != otherSession.UserID:
		return false
	case !util.CompareDates(userSession.ExpireDate, otherSession.ExpireDate):
		return false
	}

	return true
}
