package dbmodels

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserSession struct {
	Id bson.ObjectId `bson:"_id" json:"id"`

	UserId     bson.ObjectId `bson:"userId,omitempty" json:"userId"`
	Token      string        `bson:"token,omitempty" json:"token"`
	ExpireDate time.Time     `bson:"expireDate,omitempty" json:"expireDate"`
}

func (userSession *UserSession) Equal(obj Object) bool {
	otherSession, ok := obj.(*UserSession)
	if !ok {
		return false
	}

	switch {
	case userSession.Id != otherSession.Id:
		return false
	case userSession.Token != otherSession.Token:
		return false
	case userSession.UserId != otherSession.UserId:
		return false
	case !userSession.ExpireDate.Truncate(time.Millisecond).Equal(otherSession.ExpireDate.Truncate(time.Millisecond)):
		return false
	}

	return true
}
