// Package api contains the API functionality helpers
//
// Each package represents the API functionality of a
// certain endpoint which may implement some of the
// following functionalities: GET, POST, PUT or DELETE
package api

import (
	"errors"
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

var (
	// ErrEntityFormat shows that the data is not in the correct format
	ErrEntityFormat = errors.New("The entity was not in the correct format")

	// ErrEntityIntegrity shows that the data does not contain all the compulsory components
	ErrEntityIntegrity = errors.New("The entity doesn't comply to the integrity requirements")

	// ErrEntityProcess shows that the data could not be processed correctly
	ErrEntityProcess = errors.New("The entity could not be processed")

	// ErrEntityNotFound shows that the searched data was not found
	ErrEntityNotFound = errors.New("No entity with the specified data was found")

	// ErrIDParamNotSpecified shows that the ID parameter is missing from the query
	ErrIDParamNotSpecified = errors.New("No id was specified for the entity to be updated")

	// ErrInvalidIDParam shows that the ID parameter is not in the correct format
	ErrInvalidIDParam = errors.New("The userId parameter is not a valid bson.ObjectId")
)

// A Request contains the important and processable data from a HTTP request
type Request struct {
	Header        http.Header
	Form          url.Values
	ContentLength int64
	Body          []byte
	Identity      *identity.Identity
}

// A Response contains the information that will be sent back to the user
// through a HTTP response
type Response struct {
	Content      []byte
	StatusCode   int
	ErrorMessage string
	ContentType  string
	File         string
}
