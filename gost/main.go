package main

import (
	"gost/api/appuserapi"
	"gost/api/authapi"
	"gost/api/transactionapi"
	"gost/auth/cookies"
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

// APIContainer is a struct used for boxing all the existing api endpoints. It is used for mapping requests to functions.
// Add all the existing endpoints as part of this container
type APIContainer struct {
	appuserapi.ApplicationUsersAPI
	transactionapi.TransactionsAPI
	authapi.AuthAPI
}

// Application entry point - sets the behavior for the app
func main() {
	startWebFramework()
}

// startWebFramework performs the startup operations for the entire web framework
// and starts the actual http or https server used for listening for requests.
func startWebFramework() {
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
	var emptyConfigParam string

	// Initialize application configuration
	config.InitApp(emptyConfigParam)
	config.InitDatabase(emptyConfigParam)
	config.InitRoutes(emptyConfigParam)

	// Initialize the MongoDb service
	service.InitDbService()

	// Register the API endpoints
	httphandle.RegisterEndpoints(new(APIContainer))

	// Start the caching system
	cache.StartCachingSystem(cache.CacheExpireTime)

	// Initialize the cookie store in the auth module
	cookies.InitCookieStore()

	// Set the app to use all the available processors
	runtime.GOMAXPROCS(numberOfProcessors)
}

func listenForInterruptSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	log.Println("")
	log.Println("The server will now shut down gracefully...")

	service.CloseDbService()
	cache.StopCachingSystem()

	os.Exit(0)
}
