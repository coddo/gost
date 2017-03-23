package httphandle

import (
	"errors"
	"gost/api"
	"gost/auth/cookies"
	"gost/auth/identity"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"xojoc.pw/useragent"
)

var (
	// ErrNoContent shows that the underlying HTTP request doesn't contain any data/content
	ErrNoContent = errors.New("No content has been received")

	// ErrInvalidFormFormat shows that the underlying HTTP request has its data or form in a incorrect or unparsable format
	ErrInvalidFormFormat = errors.New("The request form has an invalid format")
)

func generateRequest(req *http.Request, rw http.ResponseWriter, userIdentity *identity.Identity, params httprouter.Params) *api.Request {
	statusCode, err := parseRequestContent(req)
	if err != nil {
		sendMessageResponse(statusCode, err.Error(), rw, req)
		return nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendMessageResponse(http.StatusBadRequest, err.Error(), rw, req)
		return nil
	}

	var clientDetails = parseClientDetails(req)

	request := &api.Request{
		Header:        req.Header,
		Form:          req.Form,
		RouteValues:   params,
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

// parseRequestContent performs validity checks on a request based on the HTTP method used.
// Checks are made for data content if the methods are POST or PUT, and if the url form can be correctly parsed
func parseRequestContent(request *http.Request) (int, error) {
	if request.ContentLength == 0 {
		if request.Method == "POST" || request.Method == "PUT" {
			return http.StatusBadRequest, ErrNoContent
		}
	}

	err := request.ParseForm()

	if err != nil {
		return http.StatusBadRequest, ErrInvalidFormFormat
	}

	return -1, nil
}

func respond(resp *api.Response, rw http.ResponseWriter, req *http.Request) {
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusInternalServerError
		sendMessageResponse(resp.StatusCode, api.StatusText(resp.StatusCode), rw, req)
	} else if len(resp.ErrorMessage) > 0 {
		sendMessageResponse(resp.StatusCode, resp.ErrorMessage, rw, req)
	} else {
		if len(resp.ContentType) == 0 {
			resp.ContentType = api.ContentJSON
		}

		sendResponse(resp.StatusCode, resp.Content, rw, req, resp.ContentType, resp.File)
	}
}
