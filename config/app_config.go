// Package used for application configuration
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	ENV_APPLICATION_NAME    = "GOST_TESTAPP_NAME"
	ENV_API_INSTANCE        = "GOST_TESTAPP_INSTANCE"
	ENV_HTTP_SERVER_ADDRESS = "GOST_TESTAPP_HTTP"
)

// Application configuration file path
var appConfigFile = "config/app.json"

// Application descriptive variables
var (
	ApplicationName   string
	ApiInstance       string
	HttpServerAddress string
)

// Struct with the sole purpose of easier serialization
// and deserialization of configuration data
type appConfigHolder struct {
	ApplicationName   string `json:"applicationName"`
	ApiInstance       string `json:"apiInstance"`
	HttpServerAddress string `json:"httpServerAddress"`
}

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
	ApiInstance = configData.ApiInstance
	HttpServerAddress = configData.HttpServerAddress
}

func InitTestsApp() {
	ApplicationName = os.Getenv(ENV_APPLICATION_NAME)
	ApiInstance = os.Getenv(ENV_API_INSTANCE)
	HttpServerAddress = os.Getenv(ENV_HTTP_SERVER_ADDRESS)
}
