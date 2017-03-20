package httphandle

import (
	"gost/api/dev/devapi"
	"net/http"
)

// CreateDevelopmentRoutes generates the routes that are used only in development mode
func CreateDevelopmentRoutes() {
	Routes = append(Routes,
		&Route{
			Path:           "/dev",
			Method:         http.MethodGet,
			AllowAnonymous: true,
			Action:         devapi.RouteActivateAppUser,
		},
		&Route{
			Path:           "/dev",
			Method:         http.MethodPost,
			AllowAnonymous: true,
			Action:         devapi.RouteCreateAppUser,
		},
	)
}
