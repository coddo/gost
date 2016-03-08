package config

import (
	"testing"
)

const routesFilePath = "../gost/config/routes.json"

func TestRoutesConfig(t *testing.T) {
	configRoutes(t)

	route := addRoute(t)

	r := modifyRoute(t, route.Id)

	removeRoute(t, r.Id)
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
		Id:      "TestRoute",
		Pattern: "/test/pattern/{testVar}",
		Handlers: map[string]string{
			GET_HTTP_METHOD:    "Api.GetTest",
			POST_HTTP_METHOD:   "Api.PostTest",
			PUT_HTTP_METHOD:    "Api.PutTest",
			DELETE_HTTP_METHOD: "Api.DeleteTest",
		},
	}

	err := AddRoute(&route, true)
	if err != nil {
		t.Fatal("Adding new routes failed!")
	}

	return route
}

func modifyRoute(t *testing.T, routeId string) Route {
	r := GetRoute(routeId)
	if r == nil {
		t.Fatal("Route fetching failed!")
	}

	r.Id = "TestRouteModified"
	r.Pattern = "/test/patternModified"

	err := ModifyRoute(routeId, *r, true)
	if err != nil {
		t.Fatal("Route modification failed!")
	}

	r2 := GetRoute(r.Id)
	if r2 == nil {
		t.Fatal("Modified route fetching failed!")
	}

	if !r2.Equal(*r) {
		t.Fatal("Route modification did not properly work!")
	}

	return *r2
}

func removeRoute(t *testing.T, routeId string) {
	err := RemoveRoute(routeId, true)
	if err != nil {
		t.Fatal("Removal of routes failed!")
	}

	route := GetRoute(routeId)
	if route != nil {
		t.Fatal("Route hasn't been successfully removed from the collection!")
	}
}
