package apifilter

import (
	"gost/models"
)

// CheckUserSessionIntegrity checks if an UserSession has all the compulsory fields populated
func CheckUserSessionIntegrity(userSession *models.UserSession) bool {
	switch {
	case len(userSession.ApplicationUser.ID) == 0:
		return false
	case len(userSession.Token) == 0:
		return false
	}

	return true
}
