package main

import (
	"gost/api/userapi"
	"gost/config"
	"gost/httphandle"
	"gost/servers"
	"runtime"
)

// Add all the existing endpoints as part of this container
type ApiContainer struct {
	userapi.UsersApi
}

// Function for performing automatic initializations at application startup
func initApplicationConfiguration() {
	var emptyConfigParam string = ""

	// Initialize application configuration
	config.InitApp(emptyConfigParam)
	config.InitDatabase(emptyConfigParam)
	config.InitRoutes(emptyConfigParam)

	// Initialize security module
	security.InitCrypto()

	// Register the API endpoints
	httphandle.SetApiInterface(new(ApiContainer))
}

// Application entry point - sets the behaviour for the app
func main() {
	initApplicationConfiguration()
	runtime.GOMAXPROCS(2)

	servers.StartHTTPServer()
}
