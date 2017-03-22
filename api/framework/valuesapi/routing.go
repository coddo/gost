package valuesapi

import (
	"gost/api"
)

// RouteGet performs data parsing and binding before calling the API
func RouteGet(request *api.Request) api.Response {
	return get(request.Identity)
}

// RouteGetAdmin performs data parsing and binding before calling the API
func RouteGetAdmin(request *api.Request) api.Response {
	return getAdmin(request.Identity)
}

// RouteGetAnonymous performs data parsing and binding before calling the API
func RouteGetAnonymous(request *api.Request) api.Response {
	return getAnonymous(request.Identity)
}
