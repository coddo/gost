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

	// ServerType represents the type of the server: http or https
	ServerType string

	// AccountActivationEndpoint is the link where the account activation takes place (i.e. http://example.com?token=)
	AccountActivationEndpoint string

	// PasswordResetEndpoint is the link where the password reset action takes place (i.e. http://example.com?token=)
	PasswordResetEndpoint string
)

// Struct with the sole purpose of easier serialization
// and deserialization of configuration data
type appConfigHolder struct {
	ApplicationName           string `json:"applicationName"`
	APIInstance               string `json:"apiInstance"`
	HTTPServerAddress         string `json:"httpServerAddress"`
	ServerType                string `json:"serverType"`
	AccountActivationEndpoint string `json:"accountActivationEndpoint"`
	PasswordResetEndpoint     string `json:"passwordResetEndpoint"`
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
		log.Fatalf("[InitApp] %v\n", err)
	}

	err = json.Unmarshal(data, &configData)
	if err != nil {
		log.Fatalf("[InitApp] %v\n", err)
	}

	ApplicationName = configData.ApplicationName
	APIInstance = configData.APIInstance
	HTTPServerAddress = configData.HTTPServerAddress
	ServerType = configData.ServerType
}
