package api

import (
	"gost/auth/cookies"
	"gost/auth/identity"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
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

// A Request contains the important and processable data from a HTTP request
type Request struct {
	Header        http.Header
	Form          url.Values
	RouteValues   httprouter.Params
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
