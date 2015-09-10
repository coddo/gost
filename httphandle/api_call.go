package httphandle

import (
	"gost/api"
	"gost/config"
	"gost/filter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

var apiInterface interface{}

func SetApiInterface(interf interface{}) {
	apiInterface = interf
}

func PerformApiCall(handlerName string, rw http.ResponseWriter, req *http.Request, route *config.Route) {
	// Defered function call for when the mgo driver panics
	defer func() {
		if r := recover(); r != nil {
			GiveApiStatus(StatusTooManyRequests, rw, req, route.Pattern)
		}
	}()

	// Prepare data vector for an api/endpoint call
	inputs := make([]reflect.Value, 1)

	// Create the variables containing request data
	vars := createApiVars(req, rw, route)
	if vars == nil {
		GiveApiStatus(http.StatusInternalServerError, rw, req, route.Pattern)
		return
	}

	// Populate the data vector for the api call
	inputs[0] = reflect.ValueOf(vars)

	// Perform the call on the corresponding endpoint and function
	// This is done by using reflection techniques
	var respObjects []reflect.Value
	apiMethod := reflect.ValueOf(apiInterface).MethodByName(route.Handlers[req.Method])

	// Check for zero value
	if apiMethod != *new(reflect.Value) {
		respObjects = apiMethod.Call(inputs)
	} else {
		GiveApiStatus(http.StatusInternalServerError, rw, req, route.Pattern)

		log.Println("The endpoint method is either inexistent or incorrectly mapped. Please check the server configuration files!")

		return
	}

	if respObjects == nil {
		GiveApiStatus(http.StatusInternalServerError, rw, req, route.Pattern)
		return
	}

	// Extract the response from the endpoint into a concrete type
	resp := respObjects[0].Interface().(api.ApiResponse)

	// Give the response to the api client
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusInternalServerError
		GiveApiMessage(resp.StatusCode, http.StatusText(resp.StatusCode), rw, req, route.Pattern)
	} else if len(resp.ErrorMessage) > 0 {
		GiveApiMessage(resp.StatusCode, resp.ErrorMessage, rw, req, route.Pattern)
	} else {
		if len(resp.ContentType) == 0 {
			resp.ContentType = ContentJSON
		}

		GiveApiResponse(resp.StatusCode, resp.Message, rw, req, route.Pattern, resp.ContentType, resp.File)
	}
}

func createApiVars(req *http.Request, rw http.ResponseWriter, route *config.Route) *api.ApiVar {
	err, statusCode := filter.CheckMethodAndParseContent(req)
	if err != nil {
		GiveApiMessage(statusCode, err.Error(), rw, req, route.Pattern)
		return nil
	}

	body, err := convertBodyToReadableFormat(req.Body)
	if err != nil {
		GiveApiMessage(http.StatusBadRequest, err.Error(), rw, req, route.Pattern)
		return nil
	}

	vars := &api.ApiVar{
		RequestHeader:        req.Header,
		RequestForm:          req.Form,
		RequestContentLength: req.ContentLength,
		RequestBody:          body,
	}

	return vars
}

func convertBodyToReadableFormat(data io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(data)

	return body, err
}
