package httphandle

import (
	"fmt"
	"gost/api"
	"gost/auth/cookies"
	"gost/auth/identity"
	"gost/config"
	"gost/filter"
	"io/ioutil"
	"net/http"
	"reflect"

	"xojoc.pw/useragent"
)

var endpointsContainer interface{}

var zeroEndpointMethod = *new(reflect.Value)

// RegisterEndpoints registers all the endpoints that are going to be mapped in the application
func RegisterEndpoints(container interface{}) {
	endpointsContainer = container
}

// RouteRequest parses the data from a HTTP request, determines which mapped endpoind needs to be called
// and forwards the request data to the found endpoint if it is valid.
func RouteRequest(rw http.ResponseWriter, req *http.Request, route *config.Route, endpointAction string, userIdentity *identity.Identity) {
	// Prepare recover mechanism in case of panic
	defer recoverFromError(rw, req, route.Endpoint, endpointAction)

	// Prepare data vector for an api/endpoint call
	actionParameters := make([]reflect.Value, 1)

	// Create the variables containing request data
	request := generateRequest(req, rw, route, endpointAction, userIdentity)
	if request == nil {
		return
	}

	// Populate the data vector for the api call
	actionParameters[0] = reflect.ValueOf(request)

	// Find out the name of the method where the request will be forwarded,
	// based on the registered endpoints
	endpointMethod := reflect.ValueOf(endpointsContainer).MethodByName(endpointAction)

	// Check if the searched action from the endpoint exists
	if endpointMethod == zeroEndpointMethod {
		message := "The endpoint action is either inexistent or incorrectly mapped. Please check the server configuration"
		sendMessageResponse(http.StatusInternalServerError, message, rw, req, route.Endpoint, endpointAction)
		return
	}

	// Call the mapped method from the corresponding endpoint, using the extracted and parsed data from the HTTP request
	respObjects := endpointMethod.Call(actionParameters)
	if respObjects == nil {
		sendStatusResponse(http.StatusInternalServerError, rw, req, route.Endpoint, endpointAction)
		return
	}

	// Extract the response from the endpoint into a concrete type
	resp := respObjects[0].Interface().(api.Response)

	// Give the response to the api client
	respond(&resp, rw, req, route.Endpoint, endpointAction)
}

func recoverFromError(rw http.ResponseWriter, req *http.Request, pattern, endpointAction string) {
	if err := recover(); err != nil {
		message := fmt.Sprintf("%s", err)
		sendMessageResponse(http.StatusInternalServerError, message, rw, req, pattern, endpointAction)
	}
}

func respond(resp *api.Response, rw http.ResponseWriter, req *http.Request, endpoint, endpointAction string) {
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusInternalServerError
		sendMessageResponse(resp.StatusCode, api.StatusText(resp.StatusCode), rw, req, endpoint, endpointAction)
	} else if len(resp.ErrorMessage) > 0 {
		sendMessageResponse(resp.StatusCode, resp.ErrorMessage, rw, req, endpoint, endpointAction)
	} else {
		if len(resp.ContentType) == 0 {
			resp.ContentType = api.ContentJSON
		}

		sendResponse(resp.StatusCode, resp.Content, rw, req, endpoint, endpointAction, resp.ContentType, resp.File)
	}
}

func generateRequest(req *http.Request, rw http.ResponseWriter, route *config.Route, endpointAction string, userIdentity *identity.Identity) *api.Request {
	statusCode, err := filter.ParseRequestContent(req)
	if err != nil {
		sendMessageResponse(statusCode, err.Error(), rw, req, route.Endpoint, endpointAction)
		return nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendMessageResponse(http.StatusBadRequest, err.Error(), rw, req, route.Endpoint, endpointAction)
		return nil
	}

	var clientDetails = parseClientDetails(req)

	request := &api.Request{
		Header:        req.Header,
		Form:          req.Form,
		ContentLength: req.ContentLength,
		Body:          body,
		Identity:      userIdentity,
		ClientDetails: clientDetails,
	}

	return request
}

func parseClientDetails(req *http.Request) *cookies.Client {
	var userAgent = useragent.Parse(req.UserAgent())

	if userAgent == nil {
		return cookies.UnknownClientDetails()
	}

	var client = cookies.Client{
		Address:        req.RemoteAddr,
		Type:           userAgent.Type.String(),
		Name:           userAgent.Name,
		Version:        userAgent.Version.String(),
		OS:             userAgent.OS,
		IsMobileDevice: userAgent.Mobile || userAgent.Tablet,
	}

	return &client
}
