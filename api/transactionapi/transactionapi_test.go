package transactionapi

import (
	"fmt"
	"gost/api"
	"gost/auth/identity"
	"gost/orm/models"
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

var transactionsRoute = fmt.Sprintf(`[{"id": "TransactionsRoute", "endpoint": "/transactions", 
    "actions": {"%s": "POST", "%s": "GET"}}]`, CREATE, GET)

type dummyTransaction struct {
	BadField string
}

func TestTransactionsApi(t *testing.T) {
	var id bson.ObjectId

	tests.InitializeServerConfigurations(transactionsRoute, new(TransactionsAPI))

	// Cleanup function
	defer func() {
		recover()
		transactionservice.DeleteTransaction(id)
	}()

	testPostTransactionInBadFormat(t)
	testPostTransactionNotIntegral(t)
	id = testPostTransactionInGoodFormat(t)
	testGetTransactionWithInexistentIDInDB(t)
	testGetTransactionWithBadIDParam(t)
	testGetTransactionWithGoodIDParam(t, id)
}

func testGetTransactionWithInexistentIDInDB(t *testing.T) {
	params := url.Values{}
	params.Add("transactionId", bson.NewObjectId().Hex())

	tests.PerformTestRequest(apiPath, GET, api.GET, http.StatusNotFound, params, nil, t)
}

func testGetTransactionWithBadIDParam(t *testing.T) {
	params := url.Values{}
	params.Add("transactionId", "2as456fas4")

	tests.PerformTestRequest(apiPath, GET, api.GET, http.StatusBadRequest, params, nil, t)
}

func testGetTransactionWithGoodIDParam(t *testing.T, id bson.ObjectId) {
	params := url.Values{}
	params.Add("transactionId", id.Hex())

	rw := tests.PerformTestRequest(apiPath, GET, api.GET, http.StatusOK, params, nil, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in a corrupt format:", body)
	}
}

func testPostTransactionInBadFormat(t *testing.T) {
	dTransaction := &dummyTransaction{
		BadField: "bad value",
	}

	tests.PerformTestRequest(apiPath, CREATE, api.POST, http.StatusBadRequest, nil, dTransaction, t)
}

func testPostTransactionNotIntegral(t *testing.T) {
	transaction := &models.Transaction{
		ID:       bson.NewObjectId(),
		Payer:    identity.ApplicationUser{ID: bson.NewObjectId()},
		Currency: "USD",
	}

	tests.PerformTestRequest(apiPath, CREATE, api.POST, http.StatusBadRequest, nil, transaction, t)
}

func testPostTransactionInGoodFormat(t *testing.T) bson.ObjectId {
	transaction := &models.Transaction{
		ID:       bson.NewObjectId(),
		Payer:    identity.ApplicationUser{ID: bson.NewObjectId()},
		Receiver: identity.ApplicationUser{ID: bson.NewObjectId()},
		Type:     models.TransactionTypeCash,
		Ammount:  216.365,
		Currency: "USD",
	}

	rw := tests.PerformTestRequest(apiPath, CREATE, api.POST, http.StatusCreated, nil, transaction, t)

	body := rw.Body.String()
	if len(body) == 0 {
		t.Error("Response body is empty or in deteriorated format:", body)
	}

	return transaction.ID
}
