package api

import (
    "net/http"
)

// Return a status and message that signals the API client
// about an 'internal server error' that has occured
func InternalServerError(err error) ApiResponse {
    return ApiResponse{
        StatusCode:   http.StatusInternalServerError,
        ErrorMessage: err.Error(),
    }
}

// Return a status and message that signals the API client
// about a 'bad request' that the client has made to the server
func BadRequest(err error) ApiResponse {
    return ApiResponse{
        StatusCode:   http.StatusBadRequest,
        ErrorMessage: err.Error(),
    }
}

// Return a status and message that signals the API client
// that the searched resource was not found on the server
func NotFound(err error) ApiResponse {
    return ApiResponse{
        StatusCode:   http.StatusNotFound,
        ErrorMessage: http.StatusText(http.StatusNotFound),
    }
}

// Return a status and message that signals the API client
// that the accessed endpoint is either disabled or currently
// unavailable
func ServiceUnavailable(err error) ApiResponse {
    return ApiResponse{
        StatusCode:   http.StatusServiceUnavailable,
        ErrorMessage: err.Error(),
    }
}

// Return a status and message that signals the API client
// that the used HTTP Method is not allowed on this endpoint
func MethodNotAllowed(err error) ApiResponse {
    return ApiResponse{
        StatusCode:   http.StatusMethodNotAllowed,
        ErrorMessage: err.Error(),
    }
}

// Return a status and message that signals the API client
// that the login failed or that the client isn't
// logged in and therefore not authorized to use the endpoint
func Unauthorized(err error) ApiResponse {
    return ApiResponse{
        StatusCode:   http.StatusUnauthorized,
        ErrorMessage: err.Error(),
    }
}
