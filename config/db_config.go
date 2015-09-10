package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	ENV_DB_NAME = "GOST_TESTAPP_DB_NAME"
	ENV_DB_CONN = "GOST_TESTAPP_DB_CONN"
)

// Database configuration file path
var dbConfigFileName string = "config/db.json"

// Struct for modelling the configuration json representing
// The database connection details
type DbConfig struct {
	DatabaseName             string `json:"databaseName"`
	DatabaseConnectionString string `json:"databaseConnectionString"`
}

// The database connection string variable
// This variable needs to be initialized
var DbConnectionString string

// The name of the database that will be used
// This variable needs to be initialized
var DbName string

func fetchAndDeserializeDbData(filePath string) *DbConfig {
	var configEntity DbConfig

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &configEntity)

	if err != nil {
		log.Fatal(err)
	}

	return &configEntity
}

// Initialization of production database
func InitDatabase(configFile string) {
	if len(configFile) != 0 {
		dbConfigFileName = configFile
	}

	data := fetchAndDeserializeDbData(dbConfigFileName)

	DbName = data.DatabaseName
	DbConnectionString = data.DatabaseConnectionString
}

// Initialization of tests database
func InitTestsDatabase() {
	dbName := os.Getenv(ENV_DB_NAME)
	dbConn := os.Getenv(ENV_DB_CONN)

	if len(dbName) == 0 || len(dbConn) == 0 {
		log.Fatal("Environment variables for the test database are not set!")
	}

	DbName = dbName
	DbConnectionString = dbConn
}
