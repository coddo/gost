package api

import (
	"errors"
	"net/url"
	"strconv"

	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// GetIntParameter extracts an integer value from url paramters, based on its name
func GetIntParameter(paramName string, reqForm url.Values) (int, error) {
	value := reqForm.Get(paramName)

	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal, nil
	}

	return -1, fmt.Errorf("The %s parameter is not in the correct format", paramName)
}

// GetStringParameter extracts a string value from url paramters, based on its name
func GetStringParameter(paramName string, reqForm url.Values) string {
	return reqForm.Get(paramName)
}

// GetIDParameter extracts a bson.ObjectID value from url paramters, based on its name
func GetIDParameter(paramName string, reqForm url.Values) (bson.ObjectId, error) {
	id := reqForm.Get(paramName)

	if !bson.IsObjectIdHex(id) {
		return "", errors.New("The id parameter is not a valid bson.ObjectId")
	}

	return bson.ObjectIdHex(id), nil
}

// GetIntRouteValue extracts an integer value from the url route, based on its name
func GetIntRouteValue(valueName string, routeValues map[string]string) (int, error) {
	if value, found := routeValues[valueName]; found {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal, nil
		}
	}

	return -1, fmt.Errorf("The %s route value is not in the correct format", valueName)
}

// GetStringRouteValue extracts a string value from the url route, based on its name
func GetStringRouteValue(valueName string, routeValues map[string]string) string {
	if value, found := routeValues[valueName]; found {
		return value
	}

	return ""
}

// GetIDRouteValue extracts a bson.ObjectID value from the url route, based on its name
func GetIDRouteValue(valueName string, routeValues map[string]string) (bson.ObjectId, error) {
	if value, found := routeValues[valueName]; found {
		if bson.IsObjectIdHex(value) {
			return bson.ObjectIdHex(value), nil
		}
	}

	return "", errors.New("The id parameter is not a valid bson.ObjectId")
}
