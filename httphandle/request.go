package httphandle

import (
	"gost/api"
	"gost/auth/cookies"
	"gost/auth/identity"
	"gost/filter"
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"
	"xojoc.pw/useragent"
)

func generateRequest(req *http.Request, rw http.ResponseWriter, userIdentity *identity.Identity) *api.Request {
	statusCode, err := filter.ParseRequestContent(req)
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
		RouteValues:   bone.GetAllValues(req),
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
