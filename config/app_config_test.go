package config

import (
	"testing"
)

const appFilePath = "test_data/app.json"

func TestAppConfig(t *testing.T) {
	InitApp(appFilePath)

	if len(ApplicationName) == 0 {
		t.Fatal("Application name was not properly loaded from the config file!")
	}

	if len(ApiInstance) == 0 {
		t.Fatal("Api instance version was not properly loaded from the config file!")
	}

	if len(HttpServerAddress) == 0 {
		t.Fatal("Api http server address was not properly loaded from the config file!")
	}
}
