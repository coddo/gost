package transactionapi

import (
	"gost/api"
	"gost/orm/models"
	"gost/util/jsonutil"
)

// RouteGetTransaction performs data parsing and binding before calling the API
func RouteGetTransaction(request *api.Request) api.Response {
	transactionID, err := request.GetIDRouteValue("transactionId")

	if err != nil {
		return api.BadRequest(err)
	}
	if len(transactionID) == 0 {
		return api.NotFound(err)
	}

	return getTransaction(transactionID)
}

// RouteCreateTransaction performs data parsing and binding before calling the API
func RouteCreateTransaction(request *api.Request) api.Response {
	transaction := &models.Transaction{}

	err := jsonutil.DeserializeJSON(request.Body, transaction)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return createTransaction(transaction)
}
