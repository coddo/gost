package httphandle

import (
	"gost/api/app/transactionapi"
	"gost/api/dev/devapi"
	"gost/api/framework/authapi"
	"gost/api/framework/sessionapi"
	"gost/api/framework/valuesapi"
	"net/http"
)

// CreateAPIRoutes generates the main API routes used by the application
func CreateAPIRoutes() {
	RegisterRoute("/transactions/{transactionId}", http.MethodGet, false, nil, transactionapi.RouteGetTransaction)
	RegisterRoute("/transactions", http.MethodPost, false, nil, transactionapi.RouteCreateTransaction)
}

// CreateFrameworkRoutes generates the routes used by the framework itself
func CreateFrameworkRoutes() {
	// register values api routes, used for status checking
	RegisterRoute("/values/get", http.MethodGet, false, nil, valuesapi.RouteGet)
	RegisterRoute("/values/get/admin", http.MethodGet, false, []string{"admin"}, valuesapi.RouteGetAdmin)
	RegisterRoute("/values/get/anonymous", http.MethodGet, true, nil, valuesapi.RouteGetAnonymous)

	// append the session api routes
	RegisterRoute("/session", http.MethodPost, false, nil, sessionapi.RouteCreateSession)
	RegisterRoute("/session/{userId}", http.MethodGet, false, nil, sessionapi.RouteGetAllSessions)
	RegisterRoute("/session/{token}", http.MethodDelete, false, nil, sessionapi.RouteKillSession)

	// append the authentication api routes
	RegisterRoute("/auth/activate", http.MethodPatch, false, nil, authapi.RouteActivateAccount)
	RegisterRoute("/auth/activate/{email}/resendEmail", http.MethodGet, false, nil, authapi.RouteResendAccountActivationEmail)
	RegisterRoute("/auth/password/{email}/reset", http.MethodGet, false, nil, sessionapi.RouteRequestResetPassword)
	RegisterRoute("/auth/password/reset", http.MethodPatch, false, nil, sessionapi.RouteResetPassword)
	RegisterRoute("/auth/password/change", http.MethodPatch, false, nil, sessionapi.RouteChangePassword)
}

// CreateDevelopmentRoutes generates the routes that are used only in development mode
func CreateDevelopmentRoutes() {
	RegisterRoute("/dev", http.MethodGet, true, nil, devapi.RouteActivateAppUser)
	RegisterRoute("/dev", http.MethodPost, true, nil, devapi.RouteCreateAppUser)
}
