package api

import (
	"errors"
	"gost/auth/cookies"
	"gost/auth/identity"
	"net/http"
	"net/url"
)

const (
	// GET constant represents a GET type http request
	GET = "GET"
	// POST constant represents a POST type http request
	POST = "POST"
	// PUT constant represents a PUT type http request
	PUT = "PUT"
	// DELETE constant represents a DELETE type http request
	DELETE = "DELETE"
)

const (
	// ContentTextPlain represents a HTTP transfer with simple text data
	ContentTextPlain = "text/plain"
	// ContentHTML represents a HTTP transfer with html data
	ContentHTML = "text/html"
	// ContentJSON represents a HTTP transfer with JSON data
	ContentJSON = "application/json"
)

// Common errors returned by the API
var (
	ErrEntityFormat        = errors.New("The entity was not in the correct format")
	ErrEntityIntegrity     = errors.New("The entity doesn't comply to the integrity requirements")
	ErrEntityProcessing    = errors.New("The entity could not be processed")
	ErrEntityNotFound      = errors.New("No entity with the specified data was found")
	ErrIDParamNotSpecified = errors.New("No id was specified for the entity to be updated")
	ErrInvalidIDParam      = errors.New("The id parameter is not a valid bson.ObjectId")
	ErrInvalidInput        = errors.New("The needed url paramters were inexistent or invalid")
)

// A Request contains the important and processable data from a HTTP request
type Request struct {
	Header        http.Header
	Form          url.Values
	ContentLength int64
	Body          []byte
	Identity      *identity.Identity
	ClientDetails *cookies.Client
}

// A Response contains the information that will be sent back to the user through a HTTP response
type Response struct {
	Content      []byte
	StatusCode   int
	ErrorMessage string
	ContentType  string
	File         string
}
