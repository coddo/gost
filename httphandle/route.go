package httphandle

import (
	"gost/api"
	"gost/api/app/transactionapi"
	"gost/auth/identity"
	"gost/filter"
	"gost/orm/models"
	"gost/util/jsonutil"
	"net/http"

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

var routes []*Route

func createRoutesCollection() {
	routes = []*Route{
		&Route{
			Path:           "/transactions/{transactionId}",
			Method:         http.MethodGet,
			AllowAnonymous: false,
			Roles:          []string{identity.UserRoleNormal},
			Action: func(request *api.Request) api.Response {
				transactionID, found, err := filter.GetIDParameter("transactionId", request.Form)
				if err != nil {
					return api.BadRequest(err)
				}
				if !found {
					return api.NotFound(err)
				}

				return transactionapi.GetTransaction(transactionID)
			},
		},
		&Route{
			Path:           "/transactions",
			Method:         http.MethodPost,
			AllowAnonymous: false,
			Roles:          []string{identity.UserRoleNormal},
			Action: func(request *api.Request) api.Response {
				transaction := &models.Transaction{}
				err := jsonutil.DeserializeJSON(request.Body, transaction)
				if err != nil {
					return api.BadRequest(api.ErrEntityFormat)
				}

				return transactionapi.CreateTransaction(transaction)
			},
		},
	}
}

// InitRoutes initializes the application API routes and actions
func InitRoutes(mux *bone.Mux) {
	createRoutesCollection()

	for _, route := range routes {
		mux.HandleFunc(route.Path, func(rw http.ResponseWriter, req *http.Request) {
			requestAction(rw, req, route.Method, route.AllowAnonymous, route.Roles, route.Action)
		})
	}
}
