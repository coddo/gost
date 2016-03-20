// Package config is used for application configuration management
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Application configuration file path
var appConfigFile = "config/app.json"

var (
	// ApplicationName represents the name of the application
	ApplicationName string

	// APIInstance represents the current instance (version, server, signature etc) of the api
	APIInstance string

	// HTTPServerAddress represents the address at which the HTTP server is started and listening
	HTTPServerAddress string
)

// Struct with the sole purpose of easier serialization
// and deserialization of configuration data
type appConfigHolder struct {
	ApplicationName   string `json:"applicationName"`
	APIInstance       string `json:"apiInstance"`
	HTTPServerAddress string `json:"httpServerAddress"`
}

// InitApp initializes the application by reading the functional parameters
// from the configuration file
func InitApp(appConfigPath string) {
	if len(appConfigPath) != 0 {
		appConfigFile = appConfigPath
	}

	configData := &appConfigHolder{}

	data, err := ioutil.ReadFile(appConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &configData)

	if err != nil {
		log.Fatal(err)
	}

	ApplicationName = configData.ApplicationName
	APIInstance = configData.APIInstance
	HTTPServerAddress = configData.HTTPServerAddress
}
