package config

import (
	"gost/config"
	"os"
)

const (
	envApplicationName   = "GOST_TESTAPP_NAME"
	envAPIInstance       = "GOST_TESTAPP_INSTANCE"
	envHTTPServerAddress = "GOST_TESTAPP_HTTP"
)

// InitTestsApp initializes the application used for testing
func InitTestsApp() {
	config.ApplicationName = os.Getenv(envApplicationName)
	config.APIInstance = os.Getenv(envAPIInstance)
	config.HTTPServerAddress = os.Getenv(envHTTPServerAddress)
}
