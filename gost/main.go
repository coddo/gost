package main

import (
	"gost/api/app/transactionapi"
	"gost/api/framework/authapi"
	"gost/api/framework/devapi"
	"gost/api/framework/valuesapi"
	"gost/auth/cookies"
	"gost/cache"
	"gost/config"
	"gost/httphandle"
	"gost/security"
	"gost/servers"
	"gost/service"
	"log"
	"os"
	"os/signal"
	"runtime"
)

var numberOfProcessors = runtime.NumCPU()

// FrameworkAPIContainer is a struct used for boxing the framework's api endpoints.
// Add here all the framework endpoints that should be used by your application
type FrameworkAPIContainer struct {
	authapi.AuthAPI
	valuesapi.ValuesAPI
}

// ApplicationAPIContainer is a struct used for boxing all the application's api endpoints.
// This also registers all the framework's vital endpoints
type ApplicationAPIContainer struct {
	FrameworkAPIContainer
	transactionapi.TransactionsAPI
}

// DevAPIContainer is used only for development purposes.
// Register all the necessary api endpoints in the APIContainer type as this one just inherits it
type DevAPIContainer struct {
	ApplicationAPIContainer
	devapi.DevAPI
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
	if len(os.Args) < 1 || os.Args[1] == "http" {
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
	// httphandle.RegisterEndpoints(new(ApplicationAPIContainer))   ----- Use this API container when deploying in PRODUCTION
	httphandle.RegisterEndpoints(new(DevAPIContainer)) //----- Use this API container when in development
	devapi.InitDevRoutes()                             //----- Uncomment this line when in development

	// Start the caching system
	//cache.StartCachingSystem(cache.CacheExpireTime)

	// Initialize the cookie store in the auth module
	cookies.InitCookieStore()

	// Initialize the encryption module
	security.InitCipherModule()

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
