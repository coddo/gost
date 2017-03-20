package httphandle

import (
	"gost/api"
	"gost/api/app/transactionapi"
	"gost/auth/identity"
	"net/http"
)

// Route represents an endpoint from the API
type Route struct {
	Path           string
	Method         string
	AllowAnonymous bool
	Roles          []string
	Action         func(request *api.Request) api.Response
}

// CreateAPIRoutes generates the main API routes used by the application
func CreateAPIRoutes() {
	Routes = append(Routes,
		&Route{
			Path:           "/transactions/{transactionId}",
			Method:         http.MethodGet,
			AllowAnonymous: false,
			Roles:          []string{identity.UserRoleNormal},
			Action:         transactionapi.RouteGetTransaction,
		},
		&Route{
			Path:           "/transactions",
			Method:         http.MethodPost,
			AllowAnonymous: false,
			Roles:          []string{identity.UserRoleNormal},
			Action:         transactionapi.RouteCreateTransaction,
		},
	)
}
