package api

import (
	"net/http"
)

// InternalServerError returns a status and message that signals the API client
// about an 'internal server error' that has occured
func InternalServerError(err error) Response {
	return Response{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: err.Error(),
	}
}

// BadRequest returns a status and message that signals the API client
// about a 'bad request' that the client has made to the server
func BadRequest(err error) Response {
	return Response{
		StatusCode:   http.StatusBadRequest,
		ErrorMessage: err.Error(),
	}
}

// NotFound returns a status and message that signals the API client
// that the searched resource was not found on the server
func NotFound(err error) Response {
	return Response{
		StatusCode:   http.StatusNotFound,
		ErrorMessage: StatusText(http.StatusNotFound),
	}
}

// ServiceUnavailable returns a status and message that signals the API client
// that the accessed endpoint is either disabled or currently
// unavailable
func ServiceUnavailable(err error) Response {
	return Response{
		StatusCode:   http.StatusServiceUnavailable,
		ErrorMessage: err.Error(),
	}
}

// MethodNotAllowed returns a status and message that signals the API client
// that the used HTTP Method is not allowed on this endpoint
func MethodNotAllowed() Response {
	return Response{
		StatusCode:   http.StatusMethodNotAllowed,
		ErrorMessage: StatusText(http.StatusMethodNotAllowed),
	}
}

// Unauthorized returns a status and message that signals the API client
// that the login failed or that the client isn't
// logged in and therefore not authorized to use the endpoint
func Unauthorized() Response {
	return Response{
		StatusCode:   http.StatusUnauthorized,
		ErrorMessage: StatusText(http.StatusUnauthorized),
	}
}
