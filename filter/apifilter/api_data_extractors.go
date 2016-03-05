package apifilter

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"strconv"
	"strings"
)

func GetIntValueFromParams(paramName string, reqForm url.Values) (int, error, bool) {
	value := reqForm.Get(paramName)
	if value == "" {
		errMsg := []string{"The", paramName, "parameter was not specified"}
		return -1, errors.New(strings.Join(errMsg, " ")), false
	}

	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal, nil, true
	}

	errMsg := []string{"The", paramName, "parameter is not in the correct format"}
	return -1, errors.New(strings.Join(errMsg, " ")), true
}

func GetStringValueFromParams(paramName string, reqForm url.Values) (string, bool) {
	value := reqForm.Get(paramName)

	return value, value != ""
}

func GetIdFromParams(reqForm url.Values) (bson.ObjectId, error, bool) {
	id := reqForm.Get("id")
	if id == "" {
		return "", errors.New("The id parameter was not specified"), false
	}

	if !bson.IsObjectIdHex(id) {
		return "", errors.New("The id parameter is not a valid bson.ObjectId"), true
	}

	return bson.ObjectIdHex(id), nil, true
}
