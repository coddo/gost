package apifilter

import (
	"gost/models"
)

func CheckUserIntegrity(user *models.ApplicationUser) bool {
	switch {
	case len(user.Email) == 0:
		return false
	case len(user.Password) == 0:
		return false
	}

	return true
}
