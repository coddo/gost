package transactionapi

import (
	"gost/api"
	"gost/auth/identity"
	"gost/dal/service/transactionservice"
	"net/http"

	"errors"

	"gopkg.in/mgo.v2/bson"
)

// Errors returned by the transactionapi
var (
	ErrTransactionNotFound = errors.New("There is no transaction with the specified ID")
)

// getTransaction endpoint retrieves a certain transaction based on its Id
func getTransaction(transactionID bson.ObjectId) api.Response {
	dbTransaction, err := transactionservice.GetTransaction(transactionID)
	if err != nil || dbTransaction == nil {
		return api.NotFound(ErrTransactionNotFound)
	}

	payer, _ := identity.GetUser(dbTransaction.PayerID)
	receiver, _ := identity.GetUser(dbTransaction.PayerID)
	transaction := NewFrom(dbTransaction, payer, receiver)

	return api.JSONResponse(http.StatusOK, transaction)
}

// createTransaction endpoint creates a new transaction with the valid transfer tokens and data
func createTransaction(transaction *Transaction) api.Response {
	var dbTransaction = transaction.ToDalModel()

	err := transactionservice.CreateTransaction(dbTransaction)
	if err != nil {
		return api.InternalServerError(err)
	}
	transaction.ID = dbTransaction.ID

	return api.JSONResponse(http.StatusCreated, transaction)
}
