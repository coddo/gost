package httphandle

import (
	"gost/config"
	"net/http"
	"net/url"
	"strings"
)

// RequestHandler receives, parses and validates a HTTP request, which is then routed to the corresponding endpoint
func RequestHandler(rw http.ResponseWriter, req *http.Request) {
	pattern, endpoint, parseSuccessful := parseRequestURL(req.URL)

	if !parseSuccessful {
		sendMessageResponse(http.StatusBadRequest, "The format of the request URL is invalid", rw, req, pattern)
		return
	}

	route := findRoute(pattern)

	if route == nil {
		sendMessageResponse(http.StatusNotFound, "404 - The requested page cannot be found", rw, req, pattern)
		return
	}

	if !validateEndpoint(endpoint, route) {
		sendMessageResponse(http.StatusUnauthorized, "The requested endpoint is either not implemented, or not allowed", rw, req, pattern)
		return
	}

	RouteRequest(endpoint, rw, req, route)
}

func findRoute(pattern string) *config.Route {
	for _, route := range config.Routes {
		if route.Pattern == pattern {
			return &route
		}
	}

	return nil
}

func validateEndpoint(endpoint string, route *config.Route) bool {
	if _, found := route.Handlers[endpoint]; found {
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
	endpoint := fullPath[lastSeparatorIndex+1:]

	// Get the pattern of the route
	pattern := fullPath[:lastSeparatorIndex]

	return pattern, endpoint, successfulParse
}
