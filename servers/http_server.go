package servers

import (
	"gost/config"
	"gost/httphandle"
	"log"
	"net/http"
	"time"
)

const (
	httpsCertFile = "gost.crt"
	httpsKeyFile  = "gost.key"
)

// StartHTTPServer starts a HTTP server that listens for requests
func StartHTTPServer() {
	http.HandleFunc(config.APIInstance, httphandle.RequestHandler)

	server := &http.Server{
		Addr:           config.HTTPServerAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("HTTP Server STARTED! Listening at:", config.HTTPServerAddress+config.APIInstance)
	log.Fatal(server.ListenAndServe())
}

// StartHTTPSServer starts a HTTPS server that listens for requests
func StartHTTPSServer() {
	http.HandleFunc(config.APIInstance, httphandle.RequestHandler)

	server := &http.Server{
		Addr:           config.HTTPServerAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("HTTPS Server STARTED! Listening at:", config.HTTPServerAddress+config.APIInstance)
	log.Fatal(server.ListenAndServeTLS(httpsCertFile, httpsKeyFile))
}
