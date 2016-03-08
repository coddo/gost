// Package containing the API functionality helpers
//
// Each package represents the API functionality of a
// certain endpoint which may implement some of the
// following functionalities: GET, POST, PUT or DELETE
package api

import (
	"errors"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

const (
	ContentTextPlain = "text/plain"
	ContentHTML      = "text/html"
	ContentJSON      = "application/json"
)

var (
	EntityFormatError    = errors.New("The entity was not in the correct format")
	EntityIntegrityError = errors.New("The entity doesn't comply to the integrity requirements")
	EntityProcessError   = errors.New("The entity could not be processed")
	EntityNotFoundError  = errors.New("No entity with the specified data was found")

	IdParamNotSpecifiedError = errors.New("No id was specified for the entity to be updated")
)

type ApiVar struct {
	RequestHeader        http.Header
	RequestForm          url.Values
	RequestContentLength int64
	RequestBody          []byte
}

type ApiResponse struct {
	Message      []byte
	StatusCode   int
	ErrorMessage string
	ContentType  string
	File         string
}
