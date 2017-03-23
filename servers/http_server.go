package servers

import (
	"gost/config"
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/julienschmidt/httprouter"
)

const (
	httpsCertFile = "gost.crt"
	httpsKeyFile  = "gost.key"
)

// Multiplexer represents the route handler used by the http server
var Multiplexer = httprouter.New()

// StartHTTPServer starts a HTTP server that listens for requests
func StartHTTPServer() {
	var serverAddress = fmt.Sprintf(":%s", config.HTTPServerPort)

	server := &http.Server{
		Addr:           serverAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        Multiplexer,
	}

	log.Println("HTTP Server STARTED! Listening at:", serverAddress)
	log.Fatal(server.ListenAndServe())
}

// StartHTTPSServer starts a HTTPS server that listens for requests
func StartHTTPSServer() {
	var serverAddress = fmt.Sprintf(":%s", config.HTTPServerPort)

	server := &http.Server{
		Addr:           serverAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        Multiplexer,
	}

	log.Println("HTTPS Server STARTED! Listening at:", serverAddress)
	log.Fatal(server.ListenAndServeTLS(httpsCertFile, httpsKeyFile))
}
