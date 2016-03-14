package transactionapi

import (
	"gost/api"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/transactionservice"
	"net/http"
)

// TransactionsAPI defines the API endpoint for application transactions of any kind
type TransactionsAPI int

// Get endpoint retrieves a certain transaction based on its Id
func (t *TransactionsAPI) Get(vars *api.Request) api.Response {
	transactionID, err, found := apifilter.GetIdFromParams(vars.Form)

	if found {
		if err != nil {
			return api.BadRequest(err)
		}

		dbTransaction, err := transactionservice.GetTransaction(transactionID)
		if err != nil || dbTransaction == nil {
			return api.NotFound(api.ErrEntityNotFound)
		}

		transaction := &models.Transaction{}
		transaction.Expand(dbTransaction)

		return api.SingleDataResponse(http.StatusOK, transaction)
	}

	return api.BadRequest(api.ErrIDParamNotSpecified)
}

// Create endpoint creates a new transaction with the valid transfer tokens and data
func (t *TransactionsAPI) Create(vars *api.Request) api.Response {
	transaction := &models.Transaction{}

	err := models.DeserializeJson(vars.Body, transaction)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	if !apifilter.CheckTransactionIntegrity(transaction) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	dbTransaction := transaction.Collapse()
	if dbTransaction == nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}

	err = transactionservice.CreateTransaction(dbTransaction)
	if err != nil {
		return api.InternalServerError(api.ErrEntityProcess)
	}
	transaction.Id = dbTransaction.Id

	return api.SingleDataResponse(http.StatusCreated, transaction)
}
