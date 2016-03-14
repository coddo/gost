package config

import (
	"encoding/json"
	"gost/config"
	"log"
	"os"
)

const (
	envApplicationName    = "GOST_TESTAPP_NAME"
	envAPIInstance        = "GOST_TESTAPP_INSTANCE"
	envHTTPServerAddress  = "GOST_TESTAPP_HTTP"
	envDatabaseName       = "GOST_TESTAPP_DB_NAME"
	envDatabaseConnection = "GOST_TESTAPP_DB_CONN"
)

func InitTestsApp() {
	config.ApplicationName = os.Getenv(envApplicationName)
	config.APIInstance = os.Getenv(envAPIInstance)
	config.HTTPServerAddress = os.Getenv(envHTTPServerAddress)
}

func InitTestsDatabase() {
	dbName := os.Getenv(envDatabaseName)
	dbConn := os.Getenv(envDatabaseConnection)

	if len(dbName) == 0 || len(dbConn) == 0 {
		log.Fatal("Environment variables for the test database are not set!")
	}

	config.DbName = dbName
	config.DbConnectionString = dbConn
}

func InitTestsRoutes(routesString string) {
	deserializeRoutes([]byte(routesString))
}

func deserializeRoutes(routesString []byte) {
	err := json.Unmarshal(routesString, &config.Routes)

	if err != nil {
		log.Fatal(err)
	}
}
