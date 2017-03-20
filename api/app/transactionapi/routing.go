package transactionapi

import (
	"gost/api"
	"gost/orm/models"
	"gost/util/jsonutil"
)

// RouteGetTransaction performs parameter parsing before calling the API
func RouteGetTransaction(request *api.Request) api.Response {
	transactionID, found, err := api.GetIDParameter("transactionId", request.Form)
	if err != nil {
		return api.BadRequest(err)
	}
	if !found {
		return api.NotFound(err)
	}

	return getTransaction(transactionID)
}

// RouteCreateTransaction performs parameter parsing before calling the API
func RouteCreateTransaction(request *api.Request) api.Response {
	transaction := &models.Transaction{}
	err := jsonutil.DeserializeJSON(request.Body, transaction)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return createTransaction(transaction)
}
