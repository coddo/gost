package filter

import (
	"errors"
	"net/http"
)

var (
	// ErrNoContent shows that the underlying HTTP request doesn't contain any data/content
	ErrNoContent = errors.New("No content has been received")

	// ErrInvalidFormFormat shows that the underlying HTTP request has its data or form in a incorrect or unparsable format
	ErrInvalidFormFormat = errors.New("The request form has an invalid format")
)

// CheckMethodAndParseContent performs validity checks on a request based on the HTTP method used.
// Checks are made for data content if the methods are POST or PUT, and if the url form can be correctly parsed
func ParseRequestContent(request *http.Request) (int, error) {
	if request.ContentLength == 0 {
		if request.Method == "POST" || request.Method == "PUT" {
			return http.StatusBadRequest, ErrNoContent
		}
	}

	err := request.ParseForm()

	if err != nil {
		return http.StatusBadRequest, ErrInvalidFormFormat
	}

	return -1, nil
}

// CheckNotNull verifies if the given []byte data is nil or represents the word "null".
// In the above stated cases, the method returns false
func CheckNotNull(data []byte) bool {
	if data == nil {
		return false
	}

	if string(data) == "null" {
		return false
	}

	return true
}
