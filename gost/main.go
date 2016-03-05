package main

import (
	"gost/api/userapi"
	"gost/cache"
	"gost/config"
	"gost/httphandle"
	"gost/servers"
	"gost/service"
	"log"
	"os"
	"os/signal"
	"runtime"
)

var numberOfProcessors = runtime.NumCPU()

// Add all the existing endpoints as part of this container
type ApiContainer struct {
	userapi.UsersApi
}

func StartWebFramework() {
	// Start listener for performing a graceful shutdown of the server
	go listenForInterruptSignal()

	// Start a http or and https server depending on the program arguments
	if len(os.Args) <= 1 || os.Args[1] == "http" {
		servers.StartHTTPServer()
	} else if os.Args[1] == "https" {
		servers.StartHTTPSServer()
	}
}

// Function for performing automatic initializations at application startup
func init() {
	var emptyConfigParam string = ""

	// Initialize application configuration
	config.InitApp(emptyConfigParam)
	config.InitDatabase(emptyConfigParam)
	config.InitRoutes(emptyConfigParam)

	// Initialize the MongoDb service
	service.InitDbService()

	// Register the API endpoints
	httphandle.SetApiInterface(new(ApiContainer))

	// Start the caching system
	cache.StartCachingSystem(cache.CACHE_EXPIRE_TIME)

	// Set the app to use all the available processors
	runtime.GOMAXPROCS(numberOfProcessors)
}

// Application entry point - sets the behavior for the app
func main() {
	StartWebFramework()
}

func listenForInterruptSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	log.Println("The server will now shut down gracefully...")

	service.CloseDbService()
	cache.StopCachingSystem()

	os.Exit(0)
}
