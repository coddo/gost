package models

import (
	"gost/dbmodels"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ApplicationUser represents a basic user account
type ApplicationUser struct {
	ID bson.ObjectId `json:"id"`

	Email                          string    `json:"email"`
	Password                       string    `json:"password"`
	AccountType                    int       `json:"accountType"`
	ResetPasswordToken             string    `json:"resetPasswordToken"`
	ResetPasswordTokenExpireDate   time.Time `json:"resetPasswordTokenExpireDate"`
	ActivateAccountToken           string    `json:"activateAccountToken"`
	ActivateAccountTokenExpireDate time.Time `json:"activateAccountTokenExpireDate"`
	Status                         bool      `json:"status"`
}

func (user *ApplicationUser) PopConstrains() {
	// Nothing to do here for now
}

func (user *ApplicationUser) Expand(dbUser *dbmodels.ApplicationUser) {
	user.ID = dbUser.ID
	user.Email = dbUser.Email
	user.Password = dbUser.Password
	user.AccountType = dbUser.AccountType
	user.ResetPasswordToken = dbUser.ResetPasswordToken
	user.ResetPasswordTokenExpireDate = dbUser.ResetPasswordTokenExpireDate
	user.ActivateAccountToken = dbUser.ActivateAccountToken
	user.ActivateAccountTokenExpireDate = dbUser.ActivateAccountTokenExpireDate
	user.Status = dbUser.Status

	user.PopConstrains()
}

func (user *ApplicationUser) Collapse() *dbmodels.ApplicationUser {
	dbUser := dbmodels.ApplicationUser{
		ID:                             user.ID,
		Email:                          user.Email,
		Password:                       user.Password,
		AccountType:                    user.AccountType,
		ResetPasswordToken:             user.ResetPasswordToken,
		ResetPasswordTokenExpireDate:   user.ResetPasswordTokenExpireDate,
		ActivateAccountToken:           user.ActivateAccountToken,
		ActivateAccountTokenExpireDate: user.ActivateAccountTokenExpireDate,
		Status: user.Status,
	}

	return &dbUser
}
