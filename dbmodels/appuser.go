package dbmodels

import (
	"gost/util"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// NormalUserAccountType represents a ordinary application user
	NormalUserAccountType = iota
	// AdministratorAccountType represents an administrator
	AdministratorAccountType = iota
)

const (
	// StatusAccountActivated represents an account that is active in the system
	StatusAccountActivated = true
	// StatusAccountDeactivated represents an account that is inactive in the system
	StatusAccountDeactivated = false
)

// ApplicationUser contains information necessary for managing accounts
type ApplicationUser struct {
	ID bson.ObjectId `bson:"_id" json:"id"`

	Email                          string    `bson:"email,omitempty" json:"email"`
	Password                       string    `bson:"password,omitempty" json:"password"`
	AccountType                    int       `bson:"accountType,omitempty" json:"accountType"`
	ResetPasswordToken             string    `bson:"resetPasswordToken,omitempty" json:"resetPasswordToken"`
	ResetPasswordTokenExpireDate   time.Time `bson:"resetPasswordTokenExpireDate,omitempty" json:"resetPasswordTokenExpireDate"`
	ActivateAccountToken           string    `bson:"activateAccountToken" json:"activateAccountToken"`
	ActivateAccountTokenExpireDate time.Time `bson:"activateAccountTokenExpireDate,omitempty" json:"activateAccountTokenExpireDate"`
	Status                         bool      `bson:"status,omitempty" json:"status"`
}

// Equal compares two ApplicationUser objects. Implements the Objecter interface
func (user *ApplicationUser) Equal(obj Objecter) bool {
	otherUser, ok := obj.(*ApplicationUser)
	if !ok {
		return false
	}

	switch {
	case user.ID != otherUser.ID:
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
