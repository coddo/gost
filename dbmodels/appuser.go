package dbmodels

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Account type constants
const (
	NORMAL_USER_ACCOUNT_TYPE   = 0
	ADMINISTRATOR_ACCOUNT_TYPE = 1
)

// Struct representing an user account. This is a database dbmodels
type ApplicationUser struct {
	Id bson.ObjectId `bson:"_id" json:"id"`

	Email                   string    `bson:"email,omitempty" json:"email"`
	Password                string    `bson:"password,omitempty" json:"password"`
	AccountType             int       `bson:"accountType,omitempty" json:"accountType"`
	ResetPasswordToken      string    `bson:"resetPasswordTokenToken,omitempty" json:"resetPasswordTokenToken"`
	ResetPasswordExpireDate time.Time `bson:"resetPasswordExpireDate,omitempty" json:"resetPasswordExpireDate"`
}

func (user *ApplicationUser) Equal(obj Object) bool {
	otherUser, ok := obj.(*ApplicationUser)
	if !ok {
		return false
	}

	switch {
	case user.Id != otherUser.Id:
		return false
	case user.Email != otherUser.Email:
		return false
	case user.Password != otherUser.Password:
		return false
	case user.AccountType != otherUser.AccountType:
		return false
	case user.ResetPasswordToken != otherUser.ResetPasswordToken:
		return false
	case !user.ResetPasswordExpireDate.Truncate(time.Millisecond).Equal(otherUser.ResetPasswordExpireDate.Truncate(time.Millisecond)):
		return false
	}

	return true
}
