package api

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// GetIntParameter extracts an integer value from url paramters, based on its name
func GetIntParameter(paramName string, reqForm url.Values) (int, bool, error) {
	value := reqForm.Get(paramName)
	if value == "" {
		return -1, false, nil
	}

	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal, true, nil
	}

	errMsg := []string{"The", paramName, "parameter is not in the correct format"}
	return -1, true, errors.New(strings.Join(errMsg, " "))
}

// GetStringParameter extracts a string value from url paramters, based on its name
func GetStringParameter(paramName string, reqForm url.Values) (string, bool) {
	value := reqForm.Get(paramName)

	return value, value != ""
}

// GetIDParameter extracts a bson.ObjectID value from url paramters, based on its name
func GetIDParameter(paramName string, reqForm url.Values) (bson.ObjectId, bool, error) {
	id := reqForm.Get(paramName)
	if id == "" {
		return "", false, nil
	}

	if !bson.IsObjectIdHex(id) {
		return "", true, errors.New("The id parameter is not a valid bson.ObjectId")
	}

	return bson.ObjectIdHex(id), true, nil
}
