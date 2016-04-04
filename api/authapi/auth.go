package authapi

import (
	"errors"
	"gost/api"
	"gost/auth"
	"gost/auth/cookies"
	"gost/util"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

// AuthAPI defines the API endpoint for user authorization
type AuthAPI int

type AuthModel struct {
	UserID        string          `json:"userID"`
	ClientDetails *cookies.Client `json:"clientDetails"`
}

func (a *AuthAPI) CreateSession(params *api.Request) api.Response {
	model := &AuthModel{}

	err := util.DeserializeJSON(params.Body, model)
	if err != nil {
		return api.BadRequest(err)
	}

	if !bson.IsObjectIdHex(model.UserID) {
		return api.BadRequest(errors.New("The userId parameter is not a valid bson.ObjectId"))
	}

	token, err := auth.GenerateUserAuth(bson.ObjectIdHex(model.UserID), model.ClientDetails)
	if err != nil {
		return api.BadRequest(err)
	}

	return api.ByteMsgResponse(http.StatusOK, []byte(token))
}

func (a *AuthAPI) KillSession(params *api.Request) api.Response {
	return api.BadRequest(errors.New("FUNCTIONALITY NOT IMPLEMENTED"))
}
