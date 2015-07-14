package apifilter

import (
	"go-server-template/models"
)

func CheckUserSessionIntegrity(userSession *models.UserSession) bool {
	switch {
	case len(userSession.User.Id) == 0:
		return false
	case len(userSession.Token) == 0:
		return false
	}

	return true
}
