package servers

import (
	"gost/config"
	"log"
	"net/http"
	"time"

	"github.com/go-zoo/bone"
)

const (
	httpsCertFile = "gost.crt"
	httpsKeyFile  = "gost.key"
)

// Multiplexer represents the route handler used by the http server
var Multiplexer = bone.New()

// StartHTTPServer starts a HTTP server that listens for requests
func StartHTTPServer() {
	server := &http.Server{
		Addr:           config.HTTPServerAddress + config.APIInstance,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        Multiplexer,
	}

	log.Println("HTTP Server STARTED! Listening at:", config.HTTPServerAddress+config.APIInstance)
	log.Fatal(server.ListenAndServe())
}

// StartHTTPSServer starts a HTTPS server that listens for requests
func StartHTTPSServer() {
	server := &http.Server{
		Addr:           config.HTTPServerAddress + config.APIInstance,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        Multiplexer,
	}

	log.Println("HTTPS Server STARTED! Listening at:", config.HTTPServerAddress+config.APIInstance)
	log.Fatal(server.ListenAndServeTLS(httpsCertFile, httpsKeyFile))
}
