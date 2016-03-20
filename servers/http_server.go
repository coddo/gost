package servers

import (
	"gost/config"
	"gost/httphandle"
	"log"
	"net/http"
	"time"
)

const (
	CERT_FILE = "gost.crt"
	KEY_FILE  = "gost.key"
)

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

func StartHTTPSServer() {
	http.HandleFunc(config.APIInstance, httphandle.RequestHandler)

	server := &http.Server{
		Addr:           config.HTTPServerAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("HTTPS Server STARTED! Listening at:", config.HTTPServerAddress+config.APIInstance)
	log.Fatal(server.ListenAndServeTLS(CERT_FILE, KEY_FILE))
}
