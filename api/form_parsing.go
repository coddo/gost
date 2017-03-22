package api

import (
	"errors"
	"strconv"

	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// GetIntParameter extracts an integer value from url paramters, based on its name
func (req *Request) GetIntParameter(paramName string) (int, error) {
	value := req.Form.Get(paramName)

	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal, nil
	}

	return -1, fmt.Errorf("The %s parameter is not in the correct format", paramName)
}

// GetStringParameter extracts a string value from url paramters, based on its name
func (req *Request) GetStringParameter(paramName string) string {
	return req.Form.Get(paramName)
}

// GetIDParameter extracts a bson.ObjectID value from url paramters, based on its name
func (req *Request) GetIDParameter(paramName string) (bson.ObjectId, error) {
	id := req.Form.Get(paramName)

	if !bson.IsObjectIdHex(id) {
		return "", errors.New("The id parameter is not a valid bson.ObjectId")
	}

	return bson.ObjectIdHex(id), nil
}

// GetIntRouteValue extracts an integer value from the url route, based on its name
func (req *Request) GetIntRouteValue(valueName string) (int, error) {
	if value, found := req.RouteValues[valueName]; found {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal, nil
		}
	}

	return -1, fmt.Errorf("The %s route value is not in the correct format", valueName)
}

// GetStringRouteValue extracts a string value from the url route, based on its name
func (req *Request) GetStringRouteValue(valueName string) string {
	if value, found := req.RouteValues[valueName]; found {
		return value
	}

	return ""
}

// GetIDRouteValue extracts a bson.ObjectID value from the url route, based on its name
func (req *Request) GetIDRouteValue(valueName string) (bson.ObjectId, error) {
	if value, found := req.RouteValues[valueName]; found {
		if bson.IsObjectIdHex(value) {
			return bson.ObjectIdHex(value), nil
		}
	}

	return "", errors.New("The id parameter is not a valid bson.ObjectId")
}
