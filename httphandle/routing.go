package httphandle

import (
	"gost/api"
	"net/http"

	"fmt"

	"github.com/go-zoo/bone"
)

// Route represents an endpoint from the API
type Route struct {
	Path           string
	Method         string
	AllowAnonymous bool
	Roles          []string
	Action         func(request *api.Request) api.Response
}

// Routes represents the slice of routes that are active on the server
var Routes = make([]*Route, 0)

// InitRoutes initializes the application API routes and actions
func InitRoutes(mux *bone.Mux) {
	for _, route := range Routes {
		if route.Roles == nil {
			route.Roles = make([]string, 0)
		}

		var registerFunc = getRegisterFunc(mux, route.Method)

		registerFunc(route.Path, func(rw http.ResponseWriter, req *http.Request) {
			RequestHandler(rw, req, route.Method, route.AllowAnonymous, route.Roles, route.Action)
		})
	}
}

// RequestHandler represents the main func that is called on a request once an URL match succeeds
func RequestHandler(rw http.ResponseWriter, req *http.Request, method string, allowAnonymous bool, roles []string, action func(*api.Request) api.Response) {
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
	request := generateRequest(req, rw, identity)

	// Call the endpoint
	var response = action(request)
	respond(&response, rw, req)
}

func getRegisterFunc(mux *bone.Mux, method string) func(string, http.HandlerFunc) *bone.Route {
	switch method {
	case http.MethodGet:
		return mux.GetFunc
	case http.MethodPost:
		return mux.PostFunc
	case http.MethodPut:
		return mux.PutFunc
	case http.MethodDelete:
		return mux.DeleteFunc
	case http.MethodPatch:
		return mux.PatchFunc
	case http.MethodOptions:
		return mux.OptionsFunc
	case http.MethodHead:
		return mux.HeadFunc
	}

	panic(fmt.Sprintf("HTTP Method of type %s is unsupported by the framework", method))
}
