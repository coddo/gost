package dbmodels

import (
	"gopkg.in/mgo.v2/bson"
	"gost/util"
	"time"
)

// Account type constants
const (
	NORMAL_USER_ACCOUNT_TYPE   = 0
	ADMINISTRATOR_ACCOUNT_TYPE = 1
)

// Account status constants
const (
	ACCOUNT_ACTIVATED   = true
	ACCOUNT_DEACTIVATED = false
)

// Struct representing an user account. This is a database dbmodels
type ApplicationUser struct {
	Id bson.ObjectId `bson:"_id" json:"id"`

	Email                          string    `bson:"email,omitempty" json:"email"`
	Password                       string    `bson:"password,omitempty" json:"password"`
	AccountType                    int       `bson:"accountType,omitempty" json:"accountType"`
	ResetPasswordToken             string    `bson:"resetPasswordToken,omitempty" json:"resetPasswordToken"`
	ResetPasswordTokenExpireDate   time.Time `bson:"resetPasswordTokenExpireDate,omitempty" json:"resetPasswordTokenExpireDate"`
	ActivateAccountToken           string    `bson:"activateAccountToken" json:"activateAccountToken"`
	ActivateAccountTokenExpireDate time.Time `bson:"activateAccountTokenExpireDate,omitempty" json:"activateAccountTokenExpireDate"`
	Status                         bool      `bson:"status,omitempty" json:"status"`
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
	case !util.CompareDates(user.ResetPasswordTokenExpireDate, otherUser.ResetPasswordTokenExpireDate):
		return false
	case user.ActivateAccountToken != otherUser.ActivateAccountToken:
		return false
	case !util.CompareDates(user.ActivateAccountTokenExpireDate, otherUser.ActivateAccountTokenExpireDate):
		return false
	case user.Status != otherUser.Status:
		return false
	}

	return true
}
