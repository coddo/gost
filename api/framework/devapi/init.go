package devapi

import (
	"encoding/json"
	"gost/config"
	"io/ioutil"
	"log"
)

// Routes configuration file path
var routesConfigFile = "config/devroutes.json"

// InitDevRoutes initializes the routes used for development purposes only
func InitDevRoutes() {
	routesString, err := ioutil.ReadFile(routesConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	var route = config.Route{}

	err = json.Unmarshal(routesString, &route)
	if err != nil {
		log.Fatal(err)
	}

	config.Routes = append(config.Routes, route)
}
