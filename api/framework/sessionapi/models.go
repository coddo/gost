package sessionapi

import "gopkg.in/mgo.v2/bson"

// AuthModel is a binding model used for receiving the authentication data
type LoginModel struct {
	AppUserID            bson.ObjectId `json:"appUserID"`
	Password             string        `json:"password"`
	PasswordConfirmation string        `json:"passwordConfirmation"`
}
