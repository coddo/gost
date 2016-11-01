package config

import (
	"log"
)

const (
	// Development environment
	Development = "dev"

	// Release environment
	Release = "release"
)

var envState = Development

// SetEnvironmentMode sets the type of environment the app is running in
func SetEnvironmentMode(state string) {
	envState = state

	if IsInDevMode() {
		log.Println("Application is running in development mode")
	} else {
		log.Println("Application is running in release mode")
	}
}

// IsInDevMode tells if the application is in development mode (configured for development environment)
func IsInDevMode() bool {
	switch envState {
	case Development:
		return true
	case Release:
		return false
	default:
		return true
	}
}
