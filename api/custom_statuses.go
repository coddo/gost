package api

import "net/http"

const (
	// StatusTooManyRequests is used for issuing errors regarding the number of request made by a client
	StatusTooManyRequests = 429
)

var statusText = map[int]string{
	StatusTooManyRequests: "Too many requests",
}

// StatusText returns the message associated to a http status code
func StatusText(statusCode int) string {
	msg := http.StatusText(statusCode)
	if len(msg) > 0 {
		return msg
	}

	return statusText[statusCode]
}
