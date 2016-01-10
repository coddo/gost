package dbmodels

import (
	"gopkg.in/mgo.v2/bson"
)

// Account type constants
const (
	CLIENT_ACCOUNT_TYPE        = 0
	ADMINISTRATOR_ACCOUNT_TYPE = 1
)

// Struct representing an user account. This is a database dbmodels
type User struct {
	Id bson.ObjectId `bson:"_id" json:"id"`

	Email       string `bson:"email,omitempty" json:"email"`
	Password    string `bson:"password,omitempty" json:"password"`
	AccountType int    `bson:"accountType,omitempty" json:"accountType"`
	Token       string `bson:"token,omitempty" json:"token"`

	FirstName   string `bson:"firstName,omitempty" json:"firstName"`
	MiddleName  string `bson:"middleName,omitempty" json:"middleName"`
	LastName    string `bson:"lastName,omitempty" json:"lastName"`
	CompanyName string `bson:"companyName,omitempty" json:"companyName"`
	Sex         rune   `bson:"sex,omitempty" json:"sex"`
	Country     string `bson:"country,omitempty" json:"country"`
	State       string `bson:"state,omitempty" json:"state"`
	City        string `bson:"city,omitempty" json:"city"`
	Address     string `bson:"address,omitempty" json:"address"`
	PostalCode  int    `bson:"postalCode,omitempty" json:"postalCode"`
	Picture     string `bson:"picture,omitempty" json:"picture"`
}

func (user *User) Equal(obj Object) bool {
	otherUser, ok := obj.(*User)
	if !ok {
		return false
	}

	switch {
	case user.Id != otherUser.Id:
		return false
	case user.Token != otherUser.Token:
		return false
	case user.Password != otherUser.Password:
		return false
	case user.AccountType != otherUser.AccountType:
		return false
	case user.FirstName != otherUser.FirstName:
		return false
	case user.MiddleName != otherUser.MiddleName:
		return false
	case user.LastName != otherUser.LastName:
		return false
	case user.Email != otherUser.Email:
		return false
	case user.CompanyName != otherUser.CompanyName:
		return false
	case user.Sex != otherUser.Sex:
		return false
	case user.Country != otherUser.Country:
		return false
	case user.State != otherUser.State:
		return false
	case user.City != otherUser.City:
		return false
	case user.Address != otherUser.Address:
		return false
	case user.PostalCode != otherUser.PostalCode:
		return false
	case user.Picture != otherUser.Picture:
		return false
	}

	return true
}
