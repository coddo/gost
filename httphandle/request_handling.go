package httphandle

import (
	"errors"
	"gost/api"
	"gost/auth"
	"gost/auth/identity"
	"gost/config"
	"net/http"
	"net/url"
	"strings"
)

// Authorization encompasses the identity provided by the auth package
type Authorization struct {
	Identity *identity.Identity
	Error    error
}

// RequestHandler receives, parses and validates a HTTP request, which is then routed to the corresponding endpoint and method
func RequestHandler(rw http.ResponseWriter, req *http.Request) {
	authChan := make(chan *Authorization)
	go parseAuthorizationData(req, authChan)

	endpoint, actionName, isParseSuccessful := parseRequestURL(req.URL)

	if !isParseSuccessful {
		close(authChan)
		sendMessageResponse(http.StatusBadRequest, "The format of the request URL is invalid", rw, req, endpoint, actionName)
		return
	}

	route := config.GetRoute(endpoint)

	if route == nil {
		close(authChan)
		sendMessageResponse(http.StatusNotFound, "404 - The requested page cannot be found", rw, req, endpoint, actionName)
		return
	}

	if !validateEndpoint(req.Method, actionName, route) {
		close(authChan)
		sendMessageResponse(http.StatusUnauthorized, "The requested endpoint is either not implemented, or not allowed", rw, req, endpoint, actionName)
		return
	}

	userIdentity, authError := authorize(authChan, route.Actions[actionName])
	if authError != nil {
		sendMessageResponse(http.StatusUnauthorized, authError.Error(), rw, req, route.Endpoint, actionName)
		return
	}

	RouteRequest(rw, req, route, actionName, userIdentity)
}

func authorize(authChan chan *Authorization, routeAction *config.Action) (*identity.Identity, error) {
	defer close(authChan)

	authorization := <-authChan
	if authorization.Error != nil {
		return nil, authorization.Error
	}

	user := authorization.Identity
	if (!routeAction.AllowAnonymous && user.IsAnonymous()) || (routeAction.RequireAdmin && !user.IsAdmin()) {
		return nil, errors.New(api.StatusText(http.StatusUnauthorized))
	}

	return user, nil
}

func parseAuthorizationData(req *http.Request, authChan chan *Authorization) {
	// Recover in case the authorization channel was closed before the writing is done
	defer func() {
		recover()
	}()

	identity, err := auth.Authorize(req.Header)

	authChan <- &Authorization{
		Identity: identity,
		Error:    err,
	}
}

func validateEndpoint(method, actionName string, route *config.Route) bool {
	if action, found := route.Actions[actionName]; found {
		return method == action.Type
	}

	return false
}

func parseRequestURL(u *url.URL) (string, string, bool) {
	successfulParse := true

	defer func() {
		if r := recover(); r != nil {
			successfulParse = false
		}
	}()

	fullPath := u.Path[len(config.APIInstance)-1:]

	// Cut off any URL parameters and the last '/' character if present
	if paramCharIndex := strings.Index(fullPath, "?"); paramCharIndex != -1 {
		fullPath = fullPath[:paramCharIndex]
	}

	if fullPath[len(fullPath)-1] == '/' {
		fullPath = fullPath[:len(fullPath)-1]
	}

	lastSeparatorIndex := strings.LastIndex(fullPath, "/")
	if lastSeparatorIndex == -1 {
		return "", "", false
	}

	// Get the endpoint name
	endpoint := fullPath[:lastSeparatorIndex]

	// Get the endpoint's action name
	action := fullPath[lastSeparatorIndex+1:]

	return endpoint, action, successfulParse
}
