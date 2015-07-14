package servers

import (
	"go-server-template/config"
	"go-server-template/httphandle"
	"log"
	"net/http"
	"time"
)

func StartHTTPServer() {
	http.HandleFunc(config.ApiInstance, httphandle.ApiHandler)

	server := &http.Server{
		Addr:           config.HttpServerAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("HTTP Server STARTED! Listening at:", config.HttpServerAddress)
	log.Fatal(server.ListenAndServe())
}
