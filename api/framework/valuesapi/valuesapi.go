package valuesapi

import (
	"bytes"
	"gost/api"
	"gost/auth/identity"
	"net/http"
	"strings"
)

// get performs a HTTP GET as an authorized user
func get(reqIdentity *identity.Identity) api.Response {
	var message bytes.Buffer

	message.WriteString("You are currently authorized.\nYour roles are: ")
	message.WriteString(strings.Join(reqIdentity.User.Roles, ", "))

	return api.PlainTextResponse(http.StatusOK, message.String())
}

// getAdmin performs a HTTP GET allowed only for admin users
func getAdmin(reqIdentity *identity.Identity) api.Response {
	var message bytes.Buffer

	message.WriteString("You are currently authorized as an administrator.\nYour roles are: ")
	message.WriteString(strings.Join(reqIdentity.User.Roles, ", "))

	return api.PlainTextResponse(http.StatusOK, message.String())
}

// getAnonymous performs a HTTP GET as an anonymous user
func getAnonymous(reqIdentity *identity.Identity) api.Response {
	var message bytes.Buffer
	status := http.StatusOK

	message.WriteString("You have accessed an endpoint action available for anonymous users.\n")

	if reqIdentity.IsAuthorized() {
		message.WriteString("BTW, You are an authorized user")
	} else if reqIdentity.IsAnonymous() {
		message.WriteString("BTW, You are an anonymous user")
	} else {
		message.WriteString("Cannot verify your authorization status, something is wrong")
		status = http.StatusInternalServerError
	}

	return api.PlainTextResponse(status, message.String())
}
