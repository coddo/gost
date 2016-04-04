package valuesapi

import (
	"bytes"
	"gost/api"
	"net/http"
)

// ValuesAPI defines the API endpoint for the custom user testing endpoint of any kind
type ValuesAPI int

// Get performs a HTTP GET as an authorized user
func (v *ValuesAPI) Get(params *api.Request) api.Response {
	var message bytes.Buffer

	if !params.Identity.IsAuthorized() {
		return api.Unauthorized("You must be authorized to use this endpoint action")
	}

	message.WriteString("You are currently authorized.\nYour role is:")
	if params.Identity.IsAdmin() {
		message.WriteString("ADMIN")
	} else {
		message.WriteString("NORMAL USER")
	}

	return api.PlainTextResponse(http.StatusOK, message.String())
}

// GetAnonymous performs a HTTP GET as an anonymous user
func (v *ValuesAPI) GetAnonymous(params *api.Request) api.Response {
	var message bytes.Buffer
	status := http.StatusOK

	message.WriteString("You have accessed an endpoint action available for anonymous users.\n")

	if params.Identity.IsAuthorized() {
		message.WriteString("BTW, you are an authorized user")
	} else if !params.Identity.IsAnonymous() {
		message.WriteString("Cannot verify your authorization status, something is wrong")
		status = http.StatusForbidden
	} else {
		message.WriteString("BTW, You are an anonymous user")
	}

	return api.PlainTextResponse(status, message.String())
}
