package httphandle

import (
	"gost/api/framework/valuesapi"
	"net/http"
)

// CreateFrameworkRoutes generates the routes used by the framework itself
func CreateFrameworkRoutes() {
	Routes = append(Routes,
		&Route{
			Path:           "/values/get",
			Method:         http.MethodGet,
			AllowAnonymous: false,
			Action:         valuesapi.RouteGet,
		},
		&Route{
			Path:           "/values/getAnonymous",
			Method:         http.MethodGet,
			AllowAnonymous: true,
			Action:         valuesapi.RouteGetAnonymous,
		},
	)
}
