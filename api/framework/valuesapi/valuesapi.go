package valuesapi

import (
	"bytes"
	"gost/api"
	"net/http"
	"strings"
)

// Get performs a HTTP GET as an authorized user
func Get(params *api.Request) api.Response {
	var message bytes.Buffer

	message.WriteString("You are currently authorized.\nYour roles are: ")
	message.WriteString(strings.Join(params.Identity.User.Roles, ", "))

	return api.PlainTextResponse(http.StatusOK, message.String())
}

// GetAnonymous performs a HTTP GET as an anonymous user
func GetAnonymous(params *api.Request) api.Response {
	var message bytes.Buffer
	status := http.StatusOK

	message.WriteString("You have accessed an endpoint action available for anonymous users.\n")

	if params.Identity.IsAuthorized() {
		message.WriteString("BTW, You are an authorized user")
	} else if params.Identity.IsAnonymous() {
		message.WriteString("BTW, You are an anonymous user")
	} else {
		message.WriteString("Cannot verify your authorization status, something is wrong")
		status = http.StatusInternalServerError
	}

	return api.PlainTextResponse(status, message.String())
}
