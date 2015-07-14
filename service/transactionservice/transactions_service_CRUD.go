package transactionservice

import (
	"go-server-template/dbmodels"
	"go-server-template/service"
	"gopkg.in/mgo.v2/bson"
)

const CollectionName = "transactions"

func CreateTransaction(transaction *dbmodels.Transaction) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if transaction.Id == "" {
		transaction.Id = bson.NewObjectId()
	}

	err := collection.Insert(transaction)

	return err
}

func UpdateTransaction(transaction *dbmodels.Transaction) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	if transaction.Id == "" {
		return service.NoIdSpecifiedError
	}

	err := collection.UpdateId(transaction.Id, transaction)

	return err
}

func DeleteTransaction(transactionId bson.ObjectId) error {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	err := collection.RemoveId(transactionId)

	return err
}

func GetTransaction(transactionId bson.ObjectId) (*dbmodels.Transaction, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	transaction := dbmodels.Transaction{}
	err := collection.FindId(transactionId).One(&transaction)

	return &transaction, err
}

func GetAllTransactions() ([]dbmodels.Transaction, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	var transactions []dbmodels.Transaction
	err := collection.Find(bson.M{}).All(&transactions)

	return transactions, err
}

func GetAllTransactionsLimited(limit int) ([]dbmodels.Transaction, error) {
	session, collection := service.Connect(CollectionName)
	defer session.Close()

	var transactions []dbmodels.Transaction
	err := collection.Find(bson.M{}).Limit(limit).All(&transactions)

	return transactions, err
}
