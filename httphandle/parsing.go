package httphandle

import (
	"gost/api"
	"gost/auth"
	"gost/auth/cookies"
	"gost/auth/identity"
	"gost/filter"
	"io/ioutil"
	"net/http"

	"xojoc.pw/useragent"
)

func requestAction(rw http.ResponseWriter, req *http.Request, method string, endpoint string, allowAnonymous bool, roles []string, action func(*api.Request) api.Response) {
	// Check http method
	if method != req.Method {
		sendMessageResponse(http.StatusNotFound, api.StatusText(http.StatusNotFound), rw, req, endpoint)
		return
	}

	// Try authorizing the user
	var identity, isAuthorized = authorize(req, allowAnonymous, roles)
	if !isAuthorized {
		sendMessageResponse(http.StatusUnauthorized, api.StatusText(http.StatusUnauthorized), rw, req, endpoint)
		return
	}

	// Create the request
	request := generateRequest(req, rw, endpoint, identity)

	// Call the endpoint
	var response = action(request)
	respond(&response, rw, req, endpoint)
}

func authorize(req *http.Request, allowAnonymous bool, roles []string) (*identity.Identity, bool) {
	identity, err := auth.Authorize(req.Header)
	if err != nil {
		return nil, false
	}

	if (!allowAnonymous && identity.IsAnonymous()) || !identity.HasAnyRole(roles) {
		return nil, false
	}

	return identity, true
}

func generateRequest(req *http.Request, rw http.ResponseWriter, endpoint string, userIdentity *identity.Identity) *api.Request {
	statusCode, err := filter.ParseRequestContent(req)
	if err != nil {
		sendMessageResponse(statusCode, err.Error(), rw, req, endpoint)
		return nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendMessageResponse(http.StatusBadRequest, err.Error(), rw, req, endpoint)
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

func respond(resp *api.Response, rw http.ResponseWriter, req *http.Request, endpoint string) {
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusInternalServerError
		sendMessageResponse(resp.StatusCode, api.StatusText(resp.StatusCode), rw, req, endpoint)
	} else if len(resp.ErrorMessage) > 0 {
		sendMessageResponse(resp.StatusCode, resp.ErrorMessage, rw, req, endpoint)
	} else {
		if len(resp.ContentType) == 0 {
			resp.ContentType = api.ContentJSON
		}

		sendResponse(resp.StatusCode, resp.Content, rw, req, endpoint, resp.ContentType, resp.File)
	}
}
