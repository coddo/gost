package config

import (
	"testing"
)

const routesFilePath = "../gost/config/routes.json"

func TestRoutesConfig(t *testing.T) {
	configRoutes(t)

	route := addRoute(t)

	r := modifyRoute(t, route.ID)

	removeRoute(t, r.ID)
}

func configRoutes(t *testing.T) {
	InitRoutes(routesFilePath)

	if Routes == nil {
		t.Fatal("Retrieving the routes from the configuration file has failed!")
	}

	err := SaveRoutesConfiguration()

	if err != nil {
		t.Fatal("Error while saving the routes to the configuration file!")
	}
}

func addRoute(t *testing.T) Route {
	route := Route{
		ID:      "TestRoute",
		Pattern: "/test/pattern/{testVar}",
		Handlers: map[string]string{
			GetHTTPMethod:    "Api.GetTest",
			PostHTTPMethod:   "Api.PostTest",
			PutHTTPMethod:    "Api.PutTest",
			DeleteHTTPMethod: "Api.DeleteTest",
		},
	}

	err := AddRoute(&route, true)
	if err != nil {
		t.Fatal("Adding new routes failed!")
	}

	return route
}

func modifyRoute(t *testing.T, routeID string) Route {
	r := GetRoute(routeID)
	if r == nil {
		t.Fatal("Route fetching failed!")
	}

	r.ID = "TestRouteModified"
	r.Pattern = "/test/patternModified"

	err := ModifyRoute(routeID, *r, true)
	if err != nil {
		t.Fatal("Route modification failed!")
	}

	r2 := GetRoute(r.ID)
	if r2 == nil {
		t.Fatal("Modified route fetching failed!")
	}

	if !r2.Equal(*r) {
		t.Fatal("Route modification did not properly work!")
	}

	return *r2
}

func removeRoute(t *testing.T, routeID string) {
	err := RemoveRoute(routeID, true)
	if err != nil {
		t.Fatal("Removal of routes failed!")
	}

	route := GetRoute(routeID)
	if route != nil {
		t.Fatal("Route hasn't been successfully removed from the collection!")
	}
}
