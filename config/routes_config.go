package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// Routes configuration file path
var routesConfigFile = "config/routes.json"

// Routes is a variable used for storing all the routes that the api will have
var activeRoutes []Route

// Routes returns a slice containing all the currently active routes used by the application
func Routes() []Route {
	return activeRoutes
}

// InitRoutes initializes the routes based on a configuration file
func InitRoutes(routesConfigPath string) {
	if len(routesConfigPath) != 0 {
		routesConfigFile = routesConfigPath
	}

	routesData, err := ioutil.ReadFile(routesConfigFile)

	if err != nil {
		log.Fatalf("[InitRoutes] %v\n", err)
	}

	err = json.Unmarshal(routesData, &activeRoutes)
	if err != nil {
		log.Fatalf("[InitRoutes] %v\n", err)
	}
}

// SaveRoutesConfiguration saves all the active routes (Routes slice) in json format
// into the configuration file
func SaveRoutesConfiguration() error {
	if len(activeRoutes) == 0 {
		return errors.New("There are no routes configured in order to be saved")
	}

	data, err := json.MarshalIndent(activeRoutes, "", "  ")

	if err != nil {
		return errors.New("Encoding routes slice to json failed!")
	}

	err = ioutil.WriteFile(routesConfigFile, data, os.ModeDevice)
	if err != nil {
		return err
	}

	return nil
}

// AddRoutes adds a new route and makes it active
func AddRoutes(saveChangesToConfigFile bool, newRoutes ...Route) error {
	initialLength := len(activeRoutes)

	for _, route := range newRoutes {
		existingRoute := GetRoute(route.Endpoint)
		if existingRoute != nil {
			return errors.New("Route already exists!")
		}
	}

	activeRoutes = append(activeRoutes, newRoutes...)

	err := checkCollectionModification(initialLength)

	return saveChanges(err, saveChangesToConfigFile, SaveRoutesConfiguration)
}

// RemoveRoute disables and removes a certain route
func RemoveRoute(routeID string, saveChangesToConfigFile bool) error {
	initialLength := len(activeRoutes)
	index := -1

	for ind, route := range activeRoutes {
		if route.ID == routeID {
			index = ind
			break
		}
	}

	if index == -1 {
		return errors.New("Route was not found for deletion!")
	}

	activeRoutes = append(activeRoutes[:index], activeRoutes[index+1:]...)

	err := checkCollectionModification(initialLength)

	return saveChanges(err, saveChangesToConfigFile, SaveRoutesConfiguration)
}

// ModifyRoute modifies the state and information of a certain route
func ModifyRoute(routeID string, newRouteData Route, saveChangesToConfigFile bool) error {
	for i := 0; i < len(activeRoutes); i++ {
		if activeRoutes[i].ID == routeID {
			activeRoutes[i] = newRouteData

			return saveChanges(nil, saveChangesToConfigFile, SaveRoutesConfiguration)
		}
	}

	return errors.New("Route was not found for modification!")
}

// GetRoute fetches a Route entity from the active routes list, base on its ID
func GetRoute(endpoint string) *Route {
	for _, route := range activeRoutes {
		if route.Endpoint == endpoint {
			return &route
		}
	}

	return nil
}

func checkCollectionModification(initialLength int) error {
	if initialLength == len(activeRoutes) {
		return errors.New("The route couldn't be processed for the collection!")
	}

	return nil
}

func saveChanges(err error, saveChangesToConfigFile bool, saverFunction func() error) error {
	if err != nil {
		return err
	}

	if saveChangesToConfigFile {
		return saverFunction()
	}

	return nil
}
