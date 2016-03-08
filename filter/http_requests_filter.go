package filter

import (
    "errors"
    "net/http"
)

var (
    NoContentError         = errors.New("No content has been received")
    InvalidFormFormatError = errors.New("The request form has an invalid format")
)

func CheckMethodAndParseContent(r *http.Request) (error, int) {
    if r.ContentLength == 0 {
        if r.Method == "POST" || r.Method == "PUT" {
            return NoContentError, http.StatusBadRequest
        }
    }

    err := r.ParseForm()

    if err != nil {
        return InvalidFormFormatError, http.StatusBadRequest
    }

    return nil, -1
}

func CheckNotNull(data []byte) bool {
    if data == nil {
        return false
    }

    if string(data) == "null" {
        return false
    }

    return true
}
