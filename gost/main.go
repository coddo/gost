package main

import (
	"flag"
	"gost/auth/cookies"
	"gost/config"
	"gost/dal/service"
	"gost/httphandle"
	"gost/security"
	"gost/servers"
	"log"
	"os"
	"os/signal"
	"runtime"
)

var numberOfProcessors = runtime.NumCPU()

// Application flags
var (
	envFlag = ""
)

// Application entry point - sets the behavior for the app
func main() {
	// Start listener for performing a graceful shutdown of the server
	go listenForInterruptSignal()

	// Start a http or and https server depending on the program arguments
	switch config.Protocol {
	case "http":
		servers.StartHTTPServer()
	case "https":
		servers.StartHTTPSServer()
	default:
		log.Fatalf("Unkown server type: %s", config.Protocol)
	}
}

// Function for performing automatic initializations at application startup
func init() {
	log.Println("Initializing server...")
	var emptyConfigParam string

	// Initialize application flags
	defineAppFlags()

	// Initialize application configuration
	config.SetEnvironmentMode(envFlag)
	config.InitApp(emptyConfigParam)
	config.InitDatabase(emptyConfigParam)

	// Initialize the encryption module
	security.InitCipherModule()

	// Generate the necessary routes based on environemnt
	httphandle.CreateFrameworkRoutes()
	httphandle.CreateAPIRoutes()

	if config.IsInDevMode() {
		httphandle.CreateDevelopmentRoutes()
	}

	// Initialize all the generated routes
	httphandle.InitRoutes(servers.Multiplexer)

	// Initialize the MongoDb service
	service.InitDbService()

	// Initialize the cookie store in the auth module
	cookies.InitCookieStore()

	// Set the app to use all the available processors
	runtime.GOMAXPROCS(numberOfProcessors)
}

func defineAppFlags() {
	// Define application flags
	flag.StringVar(&envFlag, "env", config.Development, "The type of environemnt in which the app is run")

	// Parse all the flags
	flag.Parse()
}

func listenForInterruptSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	log.Println("\nThe server will now shut down gracefully...")

	service.CloseDbService()

	os.Exit(0)
}
