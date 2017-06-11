package config

import "testing"

const appFilePath = "../gost/config/app.json"

func TestAppConfig(t *testing.T) {
	InitApp(appFilePath)

	if len(ApplicationName) == 0 {
		t.Fatal("Application name was not properly loaded from the config file!")
	}

	if len(APIInstance) == 0 {
		t.Fatal("Api instance version was not properly loaded from the config file!")
	}

	if len(HTTPServerAddress) == 0 {
		t.Fatal("Api http server address was not properly loaded from the config file!")
	}

	if len(AccountActivationEndpoint) == 0 {
		t.Fatal("Account activation endpoint was not properly loaded from the config file!")
	}

	if len(PasswordResetEndpoint) == 0 {
		t.Fatal("Password reset endpoint was not properly loaded from the config file!")
	}
}
