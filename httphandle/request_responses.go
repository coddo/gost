package httphandle

import (
	"gost/api"
	"gost/filter"
	"log"
	"net/http"
)

func sendResponse(statusCode int, message []byte, rw http.ResponseWriter, req *http.Request, pattern, contentType, file string) {
	// Handle redirect
	if statusCode == http.StatusTemporaryRedirect {
		http.Redirect(rw, req, string(message), statusCode)
	} else {

		// Prepend necessary headers if existent or needed
		if len(rw.Header().Get("Content-Type")) == 0 && len(contentType) > 0 {
			rw.Header().Set("Content-Type", contentType)
		}

		// Handle response type
		if len(file) > 0 {
			serveFile(rw, req, file)
		} else {
			serveRawData(statusCode, message, rw)
		}
	}

	// Log event
	logRequest(statusCode, message, req.Method, pattern)
}

func sendMessageResponse(statusCode int, message string, rw http.ResponseWriter, req *http.Request, pattern string) {
	msg := []byte(message)

	sendResponse(statusCode, msg, rw, req, pattern, api.ContentTextPlain, "")
}

func sendStatusResponse(statusCode int, rw http.ResponseWriter, req *http.Request, pattern string) string {
	msg := http.StatusText(statusCode)

	if len(msg) == 0 {
		msg = StatusText(statusCode)
	}

	sendMessageResponse(statusCode, msg, rw, req, pattern)

	return msg
}

func logRequest(statusCode int, message []byte, method, pattern string) {
	if statusCode >= 400 {
		log.Println(method, pattern, statusCode, string(message))
	} else {
		log.Println(method, pattern, statusCode)
	}
}

func serveRawData(statusCode int, message []byte, rw http.ResponseWriter) {
	if filter.CheckNotNull(message) {
		rw.WriteHeader(statusCode)
		rw.Write(message)
	} else {
		rw.WriteHeader(http.StatusNoContent)
	}
}

func serveFile(rw http.ResponseWriter, req *http.Request, file string) {
	// CORS headers
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeFile(rw, req, file)
}
