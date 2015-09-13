package main

import (
	"gost/api/userapi"
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
}

// Application entry point - sets the behaviour for the app
func main() {
	go listenForInterruptSignal()

	runtime.GOMAXPROCS(numberOfProcessors)

	// Start a http or and https server depending on the program arguments
	if len(os.Args) <= 1 || os.Args[1] == "http" {
		servers.StartHTTPServer()
	} else if os.Args[1] == "https" {
		servers.StartHTTPSServer()
	}
}

func listenForInterruptSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	log.Println("The server will now shut down gracefully...")

	service.CloseDbService()

	os.Exit(0)
}
