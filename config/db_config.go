package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

// Database configuration file path
var dbConfigFileName = "config/db.json"

// DbConfig is a struct for modelling the configuration json representing
// The database connection details
type DbConfig struct {
	DatabaseName string   `json:"databaseName"`
	User         string   `json:"user"`
	Pass         string   `json:"pass"`
	Driver       string   `json:"driver"`
	Hosts        []string `json:"hosts"`
}

// DbConnectionString represents the database connection string.
// This variable needs to be initialized
var DbConnectionString string

// DbName represents the name of the database that will be used.
// This variable needs to be initialized
var DbName string

// InitDatabase initializes the production database
func InitDatabase(configFile string) {
	if len(configFile) != 0 {
		dbConfigFileName = configFile
	}

	data := fetchAndDeserializeDbData(dbConfigFileName)

	DbName = data.DatabaseName
	DbConnectionString = createConnectionString(data)
}

func fetchAndDeserializeDbData(filePath string) DbConfig {
	var configEntity DbConfig

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatalf("[InitDatabase] %v\n", err)
	}

	err = json.Unmarshal(data, &configEntity)

	if err != nil {
		log.Fatalf("[InitDatabase] %v\n", err)
	}

	return configEntity
}

func createConnectionString(data DbConfig) string {
	var buf bytes.Buffer

	buf.WriteString(data.Driver)
	buf.WriteString("://")

	if len(data.User) > 0 {
		buf.WriteString(data.User)
		buf.WriteString(":")
		buf.WriteString(data.Pass)
		buf.WriteString("@")
	}

	nrOfHosts := len(data.Hosts)
	for i := 0; i < nrOfHosts-1; i++ {
		buf.WriteString(data.Hosts[i])
		buf.WriteString(",")
	}
	buf.WriteString(data.Hosts[nrOfHosts-1])

	buf.WriteString("/")
	buf.WriteString(data.DatabaseName)

	return buf.String()
}
