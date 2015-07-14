package main

import (
	"go-server-template/api/userapi"
	"go-server-template/config"
	"go-server-template/httphandle"
	"go-server-template/servers"
	"runtime"
)

// Add all the existing endpoints as part of this container
type ApiContainer struct {
	userapi.UsersApi
}

// Function for performing automatic initializations at application startup
func initApplicationConfiguration() {
	var emptyConfigParam string = ""

	config.InitApp(emptyConfigParam)
	config.InitDatabase(emptyConfigParam)
	config.InitRoutes(emptyConfigParam)

	// Register the API endpoints
	httphandle.SetApiInterface(new(ApiContainer))
}

// Application entry point - sets the behaviour for the app
func main() {
	initApplicationConfiguration()
	runtime.GOMAXPROCS(2)

	servers.StartHTTPServer()
}
