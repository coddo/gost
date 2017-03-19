package transactionapi

import (
	"gost/api"
	"gost/bll"
	"gost/filter/apifilter"
	"gost/orm/models"

	"gopkg.in/mgo.v2/bson"
)

// GetTransaction endpoint retrieves a certain transaction based on its Id
func GetTransaction(transactionID bson.ObjectId) api.Response {
	return bll.GetTransaction(transactionID)
}

// CreateTransaction endpoint creates a new transaction with the valid transfer tokens and data
func CreateTransaction(transaction *models.Transaction) api.Response {
	if !apifilter.CheckTransactionIntegrity(transaction) {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return bll.CreateTransaction(transaction)
}
