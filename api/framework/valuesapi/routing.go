package valuesapi

import (
	"gost/api"
)

// RouteGet performs data parsing and binding before calling the API
func RouteGet(params *api.Request) api.Response {
	return get(params.Identity)
}

// RouteGetAnonymous performs data parsing and binding before calling the API
func RouteGetAnonymous(params *api.Request) api.Response {
	return getAnonymous(params.Identity)
}
