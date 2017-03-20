package httphandle

import (
	"gost/api/app/transactionapi"
	"gost/auth/identity"
	"net/http"
)

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
