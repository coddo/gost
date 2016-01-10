package transactionapi

import (
	"gost/api"
	"gost/filter/apifilter"
	"gost/models"
	"gost/service/transactionservice"
	"net/http"
)

type TransactionsApi int

const ApiName = "transactions"

func (t *TransactionsApi) GetTransaction(vars *api.ApiVar) api.ApiResponse {
	transactionId, err, found := apifilter.GetIdFromParams(vars.RequestForm)
	if found {
		if err != nil {
			return api.BadRequest(err)
		}

		dbTransaction, err := transactionservice.GetTransaction(transactionId)
		if err != nil || dbTransaction == nil {
			return api.NotFound(api.EntityNotFoundError)
		}

		transaction := &models.Transaction{}
		transaction.Expand(dbTransaction)

		return api.SingleDataResponse(http.StatusOK, transaction)
	}

	return api.BadRequest(api.IdParamNotSpecifiedError)
}

func (t *TransactionsApi) PostTransaction(vars *api.ApiVar) api.ApiResponse {
	transaction := &models.Transaction{}

	err := models.DeserializeJson(vars.RequestBody, transaction)
	if err != nil {
		return api.BadRequest(api.EntityFormatError)
	}

	if !apifilter.CheckTransactionIntegrity(transaction) {
		return api.BadRequest(api.EntityIntegrityError)
	}

	dbTransaction := transaction.Collapse()
	if dbTransaction == nil {
		return api.InternalServerError(api.EntityProcessError)
	}

	err = transactionservice.CreateTransaction(dbTransaction)
	if err != nil {
		return api.InternalServerError(api.EntityProcessError)
	}
	transaction.Id = dbTransaction.Id

	return api.SingleDataResponse(http.StatusCreated, transaction)
}
