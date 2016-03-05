package transactionapi

import (
	"gopkg.in/mgo.v2/bson"
	"gost/api"
	"gost/dbmodels"
	"gost/models"
	"gost/service/transactionservice"
	"gost/tests"
	"net/http"
	"net/url"
	"testing"
)

const transactionsRoute = "[{\"id\": \"TransactionsRoute\", \"pattern\": \"/transactions\", \"handlers\": {\"DeleteTransaction\": \"DELETE\", \"GetTransaction\": \"GET\", \"PostTransaction\": \"POST\"}}]"
const apiPath = "/transactions"

const (
	GET  = "GetTransaction"
	POST = "PostTransaction"
)

type dummyTransaction struct {
	BadField string
}

func (transaction *dummyTransaction) PopConstrains() {}

func TestTransactionsApi(t *testing.T) {
	tests.InitializeServerConfigurations(transactionsRoute, new(TransactionsApi))

	testPostTransactionInBadFormat(t)
	testPostTransactionNotIntegral(t)
	id := testPostTransactionInGoodFormat(t)
	testGetTransactionWithInexistentIdInDB(t)
	testGetTransactionWithBadIdParam(t)
	testGetTransactionWithGoodIdParam(t, id)

	// Delete the created transaction
	transactionservice.DeleteTransaction(id)
}

func testGetTransactionWithInexistentIdInDB(t *testing.T) {
	params := url.Values{}
	params.Add("id", bson.NewObjectId().Hex())

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetTransactionWithBadIdParam(t *testing.T) {
	params := url.Values{}
	params.Add("id", "2as456fas4")

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetTransactionWithGoodIdParam(t *testing.T, id bson.ObjectId) {
	params := url.Values{}
	params.Add("id", id.Hex())

	rw := tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}
}

func testPostTransactionInBadFormat(t *testing.T) {
	dTransaction := &dummyTransaction{
		BadField: "bad value",
	}

	tests.PerformApiTestCall(apiPath, POST, api.POST, http.StatusBadRequest, nil, dTransaction, t)
}

func testPostTransactionNotIntegral(t *testing.T) {
	transaction := &models.Transaction{
		Id:       bson.NewObjectId(),
		Payer:    models.User{Id: bson.NewObjectId()},
		Currency: "USD",
	}

	tests.PerformApiTestCall(apiPath, POST, api.POST, http.StatusBadRequest, nil, transaction, t)
}

func testPostTransactionInGoodFormat(t *testing.T) bson.ObjectId {
	transaction := &models.Transaction{
		Id:       bson.NewObjectId(),
		Payer:    models.User{Id: bson.NewObjectId()},
		Receiver: models.User{Id: bson.NewObjectId()},
		Type:     dbmodels.CASH_TRANSACTION_TYPE,
		Ammount:  216.365,
		Currency: "USD",
	}

	rw := tests.PerformApiTestCall(apiPath, POST, api.POST, http.StatusCreated, nil, transaction, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return transaction.Id
}
