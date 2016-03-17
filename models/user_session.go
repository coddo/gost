package models

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"gost/service/appuserservice"
	"time"
)

type UserSession struct {
	Id bson.ObjectId `json:"id"`

	ApplicationUser ApplicationUser `json:"user"`
	Token           string          `json:"token"`
	ExpireDate      time.Time       `json:"expireDate"`
}

func (userSession *UserSession) PopConstrains() {
	dbUser, err := appuserservice.GetUser(userSession.ApplicationUser.Id)
	if err == nil {
		userSession.ApplicationUser.Expand(dbUser)
	}
}

func (userSession *UserSession) Expand(dbUserSession *dbmodels.UserSession) {
	userSession.Id = dbUserSession.ID
	userSession.ApplicationUser.Id = dbUserSession.UserID
	userSession.Token = dbUserSession.Token
	userSession.ExpireDate = dbUserSession.ExpireDate

	userSession.PopConstrains()
}

func (userSession *UserSession) Collapse() *dbmodels.UserSession {
	dbUserSession := dbmodels.UserSession{
		ID:         userSession.Id,
		UserID:     userSession.ApplicationUser.Id,
		Token:      userSession.Token,
		ExpireDate: userSession.ExpireDate,
	}

	return &dbUserSession
}
