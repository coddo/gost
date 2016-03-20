package transactionservice

import (
	"gost/dbmodels"
	"gost/service"

	"gopkg.in/mgo.v2/bson"
)

const collectionName = "transactions"

// CreateTransaction adds a new Transaction to the database
func CreateTransaction(transaction *dbmodels.Transaction) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if transaction.ID == "" {
		transaction.ID = bson.NewObjectId()
	}

	err := collection.Insert(transaction)

	return err
}

// UpdateTransaction updates an existing Transaction in the database
func UpdateTransaction(transaction *dbmodels.Transaction) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	if transaction.ID == "" {
		return service.ErrNoIDSpecified
	}

	err := collection.UpdateId(transaction.ID, transaction)

	return err
}

// DeleteTransaction removes a Transaction from the database
func DeleteTransaction(transactionID bson.ObjectId) error {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	err := collection.RemoveId(transactionID)

	return err
}

// GetTransaction retrieves an Transaction from the database, based on its ID
func GetTransaction(transactionID bson.ObjectId) (*dbmodels.Transaction, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	transaction := dbmodels.Transaction{}
	err := collection.FindId(transactionID).One(&transaction)

	return &transaction, err
}

// GetAllTransactions retrieves all the existing Transaction entities in the database
func GetAllTransactions() ([]dbmodels.Transaction, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var transactions []dbmodels.Transaction
	err := collection.Find(bson.M{}).All(&transactions)

	return transactions, err
}

// GetAllTransactionsLimited retrieves the first X Transaction entities from the database, where X is the specified limit
func GetAllTransactionsLimited(limit int) ([]dbmodels.Transaction, error) {
	session, collection := service.Connect(collectionName)
	defer session.Close()

	var transactions []dbmodels.Transaction
	err := collection.Find(bson.M{}).Limit(limit).All(&transactions)

	return transactions, err
}
