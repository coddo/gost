package apifilter

import (
	"go-server-template/models"
)

func CheckUserIntegrity(user *models.User) bool {
	switch {
	case len(user.Email) == 0:
		return false
	case len(user.Token) == 0:
		return false
	case len(user.Address) == 0:
		return false
	}

	return true
}
