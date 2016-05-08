package tests

import (
	"bytes"
	"gost/config"
	"gost/httphandle"
	"gost/orm/service"
	testconfig "gost/tests/config"
	"gost/util/jsonutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
)

// PerformTestRequest does a HTTP request with test data on a specified endpoint
func PerformTestRequest(route, endpoint, method string, expectedStatusCode int, urlParams url.Values, object interface{}, t *testing.T) *httptest.ResponseRecorder {
	generatedURL, err := generateEndpointURL(route, endpoint, urlParams)
	if err != nil {
		t.Error(err.Error())
	}

	// Serialize object that will represent the request body
	// Do nothing if no object is specified
	var jsonData []byte
	if object != nil {
		jsonData, err = jsonutil.SerializeJSON(object)

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	req, err := http.NewRequest(method, generatedURL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err.Error())
	}

	rw := httptest.NewRecorder()
	httphandle.RequestHandler(rw, req)

	if rw.Code != expectedStatusCode {
		t.Fatal("Response assertion failed! Needed status:", expectedStatusCode, "Got:", rw.Code, "Message:", rw.Body.String())
	}

	return rw
}

// InitializeServerConfigurations initializes the HTTP/HTTPS server used for unit testing
func InitializeServerConfigurations(apiInterface interface{}) {
	testconfig.InitTestsApp()

	testconfig.InitTestsDatabase()
	testconfig.InitTestsRoutes()

	service.InitDbService()

	httphandle.RegisterEndpoints(apiInterface)

	runtime.GOMAXPROCS(2)
}

func generateEndpointURL(route, endpoint string, params url.Values) (*url.URL, error) {
	buffer := &bytes.Buffer{}

	if !strings.Contains(config.HTTPServerAddress, "http://") {
		buffer.WriteString("http://")
	}

	buffer.WriteString(config.HTTPServerAddress)
	buffer.WriteString(config.APIInstance[0 : len(config.APIInstance)-1])
	buffer.WriteString(route)
	buffer.WriteRune('/')
	buffer.WriteString(endpoint)

	bufferString := buffer.String()
	bufferString = strings.Replace(bufferString, "[", "", 1)
	bufferString = strings.Replace(bufferString, "]", "", 1)

	parsedURL, err := url.Parse(bufferString)
	if parsedURL != nil && params != nil {
		parsedURL.RawQuery = params.Encode()
	}

	return parsedURL, err
}
