package models

import (
	"go-server-template/dbmodels"
	"gopkg.in/mgo.v2/bson"
)

// Struct representing an user account. This is a database dbmodels
type User struct {
	Id bson.ObjectId `json:"id"`

	Email       string `json:"email"`
	Password    string `json:"password"`
	AccountType int    `json:"accountType"`
	Token       string `json:"token"`

	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
	CompanyName string `json:"companyName"`
	Sex         rune   `json:"sex"`
	Country     string `json:"country"`
	State       string `json:"state"`
	City        string `json:"city"`
	Address     string `json:"address"`
	PostalCode  int    `json:"postalCode"`
	Picture     string `json:"picture"`
}

func (user *User) PopConstrains() {
	// Nothing to do here for now
}

func (user *User) Expand(dbUser *dbmodels.User) {
	user.Id = dbUser.Id
	user.Email = dbUser.Email
	user.Password = dbUser.Password
	user.AccountType = dbUser.AccountType
	user.Token = dbUser.Token
	user.FirstName = dbUser.FirstName
	user.MiddleName = dbUser.MiddleName
	user.LastName = dbUser.LastName
	user.CompanyName = dbUser.CompanyName
	user.Sex = dbUser.Sex
	user.Country = dbUser.Country
	user.State = dbUser.State
	user.City = dbUser.City
	user.Address = dbUser.Address
	user.PostalCode = dbUser.PostalCode
	user.Picture = dbUser.Picture

	user.PopConstrains()
}

func (user *User) Collapse() *dbmodels.User {
	dbUser := dbmodels.User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		Token:       user.Token,
		AccountType: user.AccountType,
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		CompanyName: user.CompanyName,
		Sex:         user.Sex,
		Country:     user.Country,
		State:       user.State,
		City:        user.City,
		Address:     user.Address,
		PostalCode:  user.PostalCode,
		Picture:     user.Picture,
	}

	return &dbUser
}
