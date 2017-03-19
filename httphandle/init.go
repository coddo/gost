package httphandle

import (
	"net/http"

	"github.com/go-zoo/bone"
)

// Routes represents the slice of routes that are active on the server
var Routes = make([]*Route, 0)

// InitRoutes initializes the application API routes and actions
func InitRoutes(mux *bone.Mux) {
	for _, route := range Routes {
		mux.HandleFunc(route.Path, func(rw http.ResponseWriter, req *http.Request) {
			requestAction(rw, req, route.Method, route.AllowAnonymous, route.Roles, route.Action)
		})
	}
}
