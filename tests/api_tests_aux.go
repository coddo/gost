package tests

import (
	"bytes"
	"gost/config"
	"gost/httphandle"
	"gost/models"
	"gost/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
)

func PerformApiTestCall(endpointName, method string, expectedStatusCode int, urlParams url.Values, object interface{}, t *testing.T) *httptest.ResponseRecorder {
	Url, err := generateApiUrl(endpointName, urlParams)
	if err != nil {
		t.Error(err.Error())
	}

	// Serialize object that will represent the request body
	// Do nothing if no object is specified
	var jsonData []byte
	if object != nil {
		jsonData, err = models.SerializeJson(object)

		if err != nil {
			t.Fatal(err.Error())
		}
	}

	req, err := http.NewRequest(method, Url.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err.Error())
	}

	rw := httptest.NewRecorder()
	httphandle.ApiHandler(rw, req)

	if rw.Code != expectedStatusCode {
		t.Fatal("Response assertion failed! Needed:", expectedStatusCode, "Got:", rw.Code, "Message:", rw.Body.String())
	}

	return rw
}

func InitializeServerConfigurations(routeString string, apiInterface interface{}) {
	config.InitTestsApp()
	config.InitTestsDatabase()
	config.InitTestsRoutes(routeString)

	service.InitDbService()

	httphandle.SetApiInterface(apiInterface)

	runtime.GOMAXPROCS(2)
}

func generateApiUrl(path string, params url.Values) (*url.URL, error) {
	buffer := &bytes.Buffer{}

	if !strings.Contains(config.HttpServerAddress, "http://") {
		buffer.WriteString("http://")
	}

	buffer.WriteString(config.HttpServerAddress)
	buffer.WriteString(config.ApiInstance[0 : len(config.ApiInstance)-1])
	buffer.WriteString(path)

	bufferString := buffer.String()
	bufferString = strings.Replace(bufferString, "[", "", 1)
	bufferString = strings.Replace(bufferString, "]", "", 1)

	Url, err := url.Parse(bufferString)
	if Url != nil && params != nil {
		Url.RawQuery = params.Encode()
	}

	return Url, err
}
