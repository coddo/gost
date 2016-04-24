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

// GetTransaction endpoint retrieves a certain transaction based on its Id
func (t *TransactionsAPI) GetTransaction(params *api.Request) api.Response {
	transactionID, found, err := filter.GetIDParameter("transactionId", params.Form)

	if err != nil {
		return api.BadRequest(err)
	}

	if !found {
		return api.NotFound(err)
	}

	return bll.GetTransaction(transactionID)
}

// CreateTransaction endpoint creates a new transaction with the valid transfer tokens and data
func (t *TransactionsAPI) CreateTransaction(params *api.Request) api.Response {
	transaction := &models.Transaction{}

	err := util.DeserializeJSON(params.Body, transaction)
	if err != nil || !apifilter.CheckTransactionIntegrity(transaction) {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return bll.CreateTransaction(transaction)
}
