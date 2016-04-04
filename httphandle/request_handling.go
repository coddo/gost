package httphandle

import (
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
	go authorize(req, authChan)

	endpoint, endpointAction, isParseSuccessful := parseRequestURL(req.URL)

	if !isParseSuccessful {
		close(authChan)
		sendMessageResponse(http.StatusBadRequest, "The format of the request URL is invalid", rw, req, endpoint, endpointAction)
		return
	}

	route := findRoute(endpoint)

	if route == nil {
		close(authChan)
		sendMessageResponse(http.StatusNotFound, "404 - The requested page cannot be found", rw, req, endpoint, endpointAction)
		return
	}

	if !validateEndpoint(endpointAction, route) {
		close(authChan)
		sendMessageResponse(http.StatusUnauthorized, "The requested endpoint is either not implemented, or not allowed", rw, req, endpoint, endpointAction)
		return
	}

	RouteRequest(rw, req, route, endpointAction, authChan)
}

func authorize(req *http.Request, authChan chan *Authorization) {
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

func findRoute(pattern string) *config.Route {
	for _, route := range config.Routes {
		if route.Endpoint == pattern {
			return &route
		}
	}

	return nil
}

func validateEndpoint(endpoint string, route *config.Route) bool {
	if _, found := route.Actions[endpoint]; found {
		return true
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
