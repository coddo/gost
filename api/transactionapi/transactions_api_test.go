package transactionapi

import (
	"fmt"
	"gost/api"
	"gost/dbmodels"
	"gost/models"
	"gost/service/transactionservice"
	"gost/tests"
	"net/http"
	"net/url"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

const (
	GET    = "Get"
	CREATE = "Create"
)

const apiPath = "/transactions"

var transactionsRoute = fmt.Sprintf(`[{"id": "TransactionsRoute", "pattern": "/transactions", 
    "handlers": {"%s": "POST", "%s": "GET"}}]`, CREATE, GET)

type dummyTransaction struct {
	BadField string
}

func (transaction *dummyTransaction) PopConstrains() {}

func TestTransactionsApi(t *testing.T) {
	tests.InitializeServerConfigurations(transactionsRoute, new(TransactionsAPI))

	testPostTransactionInBadFormat(t)
	testPostTransactionNotIntegral(t)
	id := testPostTransactionInGoodFormat(t)
	testGetTransactionWithInexistentIDInDB(t)
	testGetTransactionWithBadIDParam(t)
	testGetTransactionWithGoodIDParam(t, id)

	// Delete the created transaction
	transactionservice.DeleteTransaction(id)
}

func testGetTransactionWithInexistentIDInDB(t *testing.T) {
	params := url.Values{}
	params.Add("id", bson.NewObjectId().Hex())

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetTransactionWithBadIDParam(t *testing.T) {
	params := url.Values{}
	params.Add("id", "2as456fas4")

	tests.PerformApiTestCall(apiPath, GET, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetTransactionWithGoodIDParam(t *testing.T, id bson.ObjectId) {
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

	tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusBadRequest, nil, dTransaction, t)
}

func testPostTransactionNotIntegral(t *testing.T) {
	transaction := &models.Transaction{
		Id:       bson.NewObjectId(),
		Payer:    models.ApplicationUser{Id: bson.NewObjectId()},
		Currency: "USD",
	}

	tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusBadRequest, nil, transaction, t)
}

func testPostTransactionInGoodFormat(t *testing.T) bson.ObjectId {
	transaction := &models.Transaction{
		Id:       bson.NewObjectId(),
		Payer:    models.ApplicationUser{Id: bson.NewObjectId()},
		Receiver: models.ApplicationUser{Id: bson.NewObjectId()},
		Type:     dbmodels.CashTransactionType,
		Ammount:  216.365,
		Currency: "USD",
	}

	rw := tests.PerformApiTestCall(apiPath, CREATE, api.POST, http.StatusCreated, nil, transaction, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return transaction.Id
}
