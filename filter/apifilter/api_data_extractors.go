package apifilter

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// GetIntValueFromParams extracts an integer value from url paramters, based on its name
func GetIntValueFromParams(paramName string, reqForm url.Values) (int, bool, error) {
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

// GetStringValueFromParams extracts a string value from url paramters, based on its name
func GetStringValueFromParams(paramName string, reqForm url.Values) (string, bool) {
	value := reqForm.Get(paramName)

	return value, value != ""
}

// GetIDFromParams extracts a bson.ObjectID value from url paramters, based on its name
func GetIDFromParams(reqForm url.Values) (bson.ObjectId, bool, error) {
	id := reqForm.Get("id")
	if id == "" {
		return "", false, nil
	}

	if !bson.IsObjectIdHex(id) {
		return "", true, errors.New("The id parameter is not a valid bson.ObjectId")
	}

	return bson.ObjectIdHex(id), true, nil
}
