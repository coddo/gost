package models

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"gost/service/userservice"
	"time"
)

type UserSession struct {
	Id bson.ObjectId `json:"id"`

	User       User      `json:"user"`
	Token      string    `json:"token"`
	ExpireDate time.Time `json:"expireDate"`
}

func (userSession *UserSession) PopConstrains() {
	dbUser, err := userservice.GetUser(userSession.User.Id)
	if err == nil {
		userSession.User.Expand(dbUser)
	}
}

func (userSession *UserSession) Expand(dbUserSession *dbmodels.UserSession) {
	userSession.Id = dbUserSession.Id
	userSession.User.Id = dbUserSession.UserId
	userSession.Token = dbUserSession.Token
	userSession.ExpireDate = dbUserSession.ExpireDate

	userSession.PopConstrains()
}

func (userSession *UserSession) Collapse() *dbmodels.UserSession {
	dbUserSession := dbmodels.UserSession{
		Id:         userSession.Id,
		UserId:     userSession.User.Id,
		Token:      userSession.Token,
		ExpireDate: userSession.ExpireDate,
	}

	return &dbUserSession
}
