package httphandle

import (
	"gost/api"
	"gost/cache"
	"gost/config"
	"gost/filter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

var endpointsContainer interface{}

// RegisterEndpoints registers all the endpoints that are going to be mapped in the application
func RegisterEndpoints(container interface{}) {
	endpointsContainer = container
}

// PerformAPICall parses the data from a HTTP request, determines which mapped endpoind needs to be called
// and forwards the request data to the found endpoint if it is valid.
func PerformAPICall(endpoint string, rw http.ResponseWriter, req *http.Request, route *config.Route) {
	// Prepare data vector for an api/endpoint call
	inputs := make([]reflect.Value, 1)

	// Create the variables containing request data
	vars := createAPIVars(req, rw, route)
	if vars == nil {
		return
	}

	// Try giving the response directly from the cache if available or invalidate it if necessary
	if respondFromCache(rw, req, route) {
		return
	}

	// Populate the data vector for the api call
	inputs[0] = reflect.ValueOf(vars)

	// Find out the name of the method where the request will be forwarded,
	// based on the registered endpoints
	apiMethod := reflect.ValueOf(endpointsContainer).MethodByName(endpoint)

	// Check if the searched method from the endpoint exists
	if apiMethod == *new(reflect.Value) {
		log.Println("The endpoint method is either inexistent or incorrectly mapped. Please check the server configuration files!")

		GiveApiStatus(http.StatusInternalServerError, rw, req, route.Pattern)

		return
	}

	// Call the mapped method from the corresponding endpoint, using the extracted and parsed data from the HTTP request
	respObjects := apiMethod.Call(inputs)
	if respObjects == nil {
		GiveApiStatus(http.StatusInternalServerError, rw, req, route.Pattern)
		return
	}

	// Extract the response from the endpoint into a concrete type
	resp := respObjects[0].Interface().(api.Response)

	// Give the response to the api client
	respond(&resp, rw, req, route.Pattern)
}

func respondFromCache(rw http.ResponseWriter, req *http.Request, route *config.Route) bool {
	if cache.Status == cache.StatusOFF {
		return false
	}

	if cachedData, err := cache.QueryByRequest(route.Pattern); err == nil {
		if req.Method == api.GET {
			GiveApiResponse(cachedData.StatusCode, cachedData.Data, rw, req, route.Pattern, cachedData.ContentType, cachedData.File)
			return true
		}

		// Invalidate the cache if a modification, deletion or addition was made to this endpoint
		cachedData.Invalidate()
	}

	return false
}

func respond(resp *api.Response, rw http.ResponseWriter, req *http.Request, endpoint string) {
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusInternalServerError
		GiveApiMessage(resp.StatusCode, http.StatusText(resp.StatusCode), rw, req, endpoint)
	} else if len(resp.ErrorMessage) > 0 {
		GiveApiMessage(resp.StatusCode, resp.ErrorMessage, rw, req, endpoint)
	} else {
		if len(resp.ContentType) == 0 {
			resp.ContentType = CONTENT_JSON
		}

		GiveApiResponse(resp.StatusCode, resp.Message, rw, req, endpoint, resp.ContentType, resp.File)

		// Try caching the data only if a GET request was made
		go func(resp *api.Response, req *http.Request, endpoint string) {
			if req.Method == api.GET && cache.Status == cache.StatusON {
				cacheResponse(resp, endpoint)
			}
		}(resp, req, endpoint)
	}
}

func cacheResponse(resp *api.Response, endpoint string) {
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) || len(resp.Message) == 0 {
		return
	}

	cacheEntity := &cache.Cache{
		Key:         cache.MapKey(endpoint),
		Data:        resp.Message,
		StatusCode:  resp.StatusCode,
		ContentType: resp.ContentType,
		File:        resp.File,
	}

	cacheEntity.Cache()
}

func createAPIVars(req *http.Request, rw http.ResponseWriter, route *config.Route) *api.Request {
	statusCode, err := filter.CheckMethodAndParseContent(req)
	if err != nil {
		GiveApiMessage(statusCode, err.Error(), rw, req, route.Pattern)
		return nil
	}

	body, err := convertBodyToReadableFormat(req.Body)
	if err != nil {
		GiveApiMessage(http.StatusBadRequest, err.Error(), rw, req, route.Pattern)
		return nil
	}

	vars := &api.Request{
		Header:        req.Header,
		Form:          req.Form,
		ContentLength: req.ContentLength,
		Body:          body,
	}

	return vars
}

func convertBodyToReadableFormat(data io.ReadCloser) ([]byte, error) {
	return ioutil.ReadAll(data)
}
