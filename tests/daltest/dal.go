package daltest

import (
	"gost/config"
	"log"
	"os"
)

const (
	envDatabaseName       = "GOST_TESTAPP_DB_NAME"
	envDatabaseConnection = "GOST_TESTAPP_DB_CONN"
)

// InitTestsDatabase initializes the connection parameters to the database used for testing
func InitTestsDatabase() {
	dbName := os.Getenv(envDatabaseName)
	dbConn := os.Getenv(envDatabaseConnection)

	if len(dbName) == 0 || len(dbConn) == 0 {
		log.Fatal("Environment variables for the test database are not set!")
	}

	config.DbName = dbName
	config.DbConnectionString = dbConn
}
