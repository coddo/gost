package models

import (
	"encoding/json"
	"go-server-template/dbmodels"
	"go-server-template/service/userservice"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserSession struct {
	Id bson.ObjectId `json:"id"`

	User       User      `json:"user"`
	Token      string    `json:"token"`
	ExpireDate time.Time `json:"expireDate"`
}

func (userSession *UserSession) SerializeJson() ([]byte, error) {
	data, err := json.MarshalIndent(*userSession, JsonPrefix, JsonIndent)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (userSession *UserSession) DeserializeJson(obj []byte) error {
	err := json.Unmarshal(obj, userSession)

	if err != nil {
		return err
	}

	return nil
}

func (userSession *UserSession) PopConstraints() {
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

	userSession.PopConstraints()
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
