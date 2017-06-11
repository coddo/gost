package httphandle

import (
	"gost/api"
	"gost/filter"
	"log"
	"net/http"
)

func sendResponse(statusCode int, message []byte, rw http.ResponseWriter, req *http.Request, contentType, filePath string) {
	// Handle redirect
	if statusCode == http.StatusTemporaryRedirect {
		http.Redirect(rw, req, string(message), statusCode)
	} else {

		// Prepend necessary headers if existent or needed
		if len(rw.Header().Get("Content-Type")) == 0 && len(contentType) > 0 {
			rw.Header().Set("Content-Type", contentType)
		}

		// Handle response type
		if len(filePath) > 0 {
			serveFile(rw, req, filePath)
		} else {
			serveRawData(statusCode, message, rw)
		}
	}

	// Log event
	go logRequest(statusCode, message, req)
}

func sendMessageResponse(statusCode int, message string, rw http.ResponseWriter, req *http.Request) {
	msg := []byte(message)

	sendResponse(statusCode, msg, rw, req, api.ContentTextPlain, "")
}

func sendStatusResponse(statusCode int, rw http.ResponseWriter, req *http.Request) string {
	message := api.StatusText(statusCode)

	sendMessageResponse(statusCode, message, rw, req)

	return message
}

func logRequest(statusCode int, message []byte, req *http.Request) {
	if statusCode >= 400 {
		log.Println(req.Method, req.RequestURI, statusCode, string(message))
	} else {
		log.Println(req.Method, req.RequestURI, statusCode)
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
