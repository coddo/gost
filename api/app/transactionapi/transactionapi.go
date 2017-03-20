package transactionapi

import (
	"gost/api"
	"gost/filter/apifilter"
	"gost/orm/models"
	"net/http"
	"testapp/service/transactionservice"

	"gopkg.in/mgo.v2/bson"
)

// GetTransaction endpoint retrieves a certain transaction based on its Id
func GetTransaction(transactionID bson.ObjectId) api.Response {
	dbTransaction, err := transactionservice.GetTransaction(transactionID)
	if err != nil || dbTransaction == nil {
		return api.NotFound(api.ErrEntityNotFound)
	}

	transaction := &models.Transaction{}
	transaction.Expand(dbTransaction)

	return api.JSONResponse(http.StatusOK, transaction)
}

// CreateTransaction endpoint creates a new transaction with the valid transfer tokens and data
func CreateTransaction(transaction *models.Transaction) api.Response {
	if !apifilter.CheckTransactionIntegrity(transaction) {
		return api.BadRequest(api.ErrEntityIntegrity)
	}

	dbTransaction := transaction.Collapse()
	if dbTransaction == nil {
		return api.InternalServerError(api.ErrEntityProcessing)
	}

	err := transactionservice.CreateTransaction(dbTransaction)
	if err != nil {
		return api.InternalServerError(api.ErrEntityProcessing)
	}
	transaction.ID = dbTransaction.ID

	return api.JSONResponse(http.StatusCreated, transaction)
}
