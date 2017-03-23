package httphandle

import (
	"gost/api"
	"net/http"

	"strings"

	"github.com/julienschmidt/httprouter"
)

// Router is the type that all the api endpoints' entry points should have
type Router func(request *api.Request) api.Response

// Route represents an endpoint from the API
type Route struct {
	Path           string
	Method         string
	AllowAnonymous bool
	Roles          []string
	Action         Router
}

// Routes represents the slice of routes that are active on the server
var routes = make([]*Route, 0)

// InitRoutes initializes the application API routes and actions
func InitRoutes(router *httprouter.Router) {
	for i := 0; i < len(routes); i++ {
		var route = routes[i]

		if route.Roles == nil {
			route.Roles = make([]string, 0)
		}

		router.Handle(route.Method, route.Path, func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
			RequestHandler(rw, req, route.Method, route.AllowAnonymous, route.Roles, route.Action, params)
		})
	}
}

// RegisterRoute registers a new route in the system. The url query will be stripped from the path, as it is used only for readability
func RegisterRoute(path, method string, allowAnonymous bool, roles []string, routeAction Router) {
	if roles == nil {
		roles = make([]string, 0)
	}

	path = strings.Split(path, "?")[0]

	routes = append(routes, &Route{
		Path:           path,
		Method:         method,
		AllowAnonymous: allowAnonymous,
		Roles:          roles,
		Action:         routeAction,
	})
}

// RequestHandler represents the main func that is called on a request once an URL match succeeds
func RequestHandler(rw http.ResponseWriter, req *http.Request, method string, allowAnonymous bool, roles []string, action Router, params httprouter.Params) {
	// Check http method
	if method != req.Method {
		sendMessageResponse(http.StatusNotFound, api.StatusText(http.StatusNotFound), rw, req)
		return
	}

	// Try authorizing the user
	var identity, isAuthorized, errorStatusCode = authorize(req, allowAnonymous, roles)
	if !isAuthorized {
		sendMessageResponse(errorStatusCode, api.StatusText(errorStatusCode), rw, req)
		return
	}

	// Create the request
	request := generateRequest(req, rw, identity, params)

	// Call the endpoint
	var response = action(request)
	respond(&response, rw, req)
}
