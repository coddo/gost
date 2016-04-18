package transactionapi

import (
	"gost/api"
	"gost/bll"
	"gost/filter"
	"gost/filter/apifilter"
	"gost/orm/models"
	"gost/util"
)

// TransactionsAPI defines the API endpoint for application transactions of any kind
type TransactionsAPI int

// Get endpoint retrieves a certain transaction based on its Id
func (t *TransactionsAPI) Get(params *api.Request) api.Response {
	transactionID, found, err := filter.GetIDFromParams("transactionId", params.Form)

	if err != nil {
		return api.BadRequest(err)
	}

	if !found {
		return api.NotFound(err)
	}

	return bll.GetTransaction(transactionID)
}

// Create endpoint creates a new transaction with the valid transfer tokens and data
func (t *TransactionsAPI) Create(params *api.Request) api.Response {
	transaction := &models.Transaction{}

	err := util.DeserializeJSON(params.Body, transaction)
	if err != nil || !apifilter.CheckTransactionIntegrity(transaction) {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return bll.CreateTransaction(transaction)
}
