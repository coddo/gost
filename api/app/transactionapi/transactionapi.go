package transactionapi

import (
	"gost/api"
	"gost/orm/models"
	"gost/orm/service/transactionservice"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

// getTransaction endpoint retrieves a certain transaction based on its Id
func getTransaction(transactionID bson.ObjectId) api.Response {
	dbTransaction, err := transactionservice.GetTransaction(transactionID)
	if err != nil || dbTransaction == nil {
		return api.NotFound(api.ErrEntityNotFound)
	}

	transaction := &models.Transaction{}
	transaction.Expand(dbTransaction)

	return api.JSONResponse(http.StatusOK, transaction)
}

// createTransaction endpoint creates a new transaction with the valid transfer tokens and data
func createTransaction(transaction *models.Transaction) api.Response {
	var dbTransaction = transaction.Collapse()
	if dbTransaction == nil {
		return api.InternalServerError(api.ErrEntityProcessing)
	}

	err := transactionservice.CreateTransaction(dbTransaction)
	if err != nil {
		return api.InternalServerError(err)
	}
	transaction.ID = dbTransaction.ID

	return api.JSONResponse(http.StatusCreated, transaction)
}
