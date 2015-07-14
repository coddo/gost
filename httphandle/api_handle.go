package httphandle

import (
	"go-server-template/config"
	"net/http"
	"net/url"
)

func ApiHandler(rw http.ResponseWriter, req *http.Request) {
	path, err := parseRequestURL(req.URL)

	if err != nil {
		GiveApiMessage(http.StatusBadRequest, "The format of the request URL is invalid", rw, req, path)
		return
	}

	route := findRoute(path)

	if route == nil {
		GiveApiMessage(http.StatusNotFound, "The requested URL cannot be found", rw, req, path)
		return
	}

	handler := findApiMethod(req.Method, route)

	if handler == "" {
		GiveApiMessage(http.StatusBadRequest, "The requested method is either not implemented, or not allowed", rw, req, path)
		return
	}

	PerformApiCall(handler, rw, req, route)
}

func findRoute(pattern string) *config.Route {
	for _, route := range config.Routes {
		if route.Pattern == pattern {
			return &route
		}
	}

	return nil
}

func findApiMethod(requestMethod string, route *config.Route) string {
	if handler, found := route.Handlers[requestMethod]; found {
		return handler
	}

	return ""
}

func parseRequestURL(u *url.URL) (string, error) {
	return u.Path[len(config.ApiInstance)-1:], nil
}
