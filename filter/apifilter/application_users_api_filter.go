package apifilter

import (
	"gost/orm/models"
)

// CheckUserIntegrity checks if an ApplicationUser has all the compulsory fields populated
func CheckUserIntegrity(user *models.ApplicationUser) bool {
	switch {
	case len(user.Email) == 0:
		return false
	case len(user.Password) == 0:
		return false
	}

	return true
}
