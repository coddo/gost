package authapi

import (
	"errors"
	"gost/api"
)

// UserSessionsAPI defines the API endpoint for user sessions
type AuthAPI int

var (
	// ErrTokenNotSpecified tells the user the token parameter was not specified
	ErrTokenNotSpecified = errors.New("The session token hasn't been specified")

	// ErrTokenExpired tells the user that the given token has expired
	ErrTokenExpired = errors.New("The session with the specified token has expired")
)

const (
	daysUntilExpire = 7
)

// Get endpoint retrieves a user session entity based on a given token
func (userSessionsApi *AuthAPI) Get(vars *api.Request) api.Response {
	// token, found := apifilter.GetStringValueFromParams("token", vars.Form)

	// if !found {
	// 	return api.BadRequest(ErrTokenNotSpecified)
	// }

	return api.MethodNotAllowed()
}

// Create endpoint creates a new user session
func (userSessionsApi *AuthAPI) Create(vars *api.Request) api.Response {
	return api.MethodNotAllowed()
}
