package models

import (
	"encoding/json"
	"go-server-template/dbmodels"
	"go-server-template/service/userservice"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Transaction struct {
	Id bson.ObjectId `json:"id"`

	Payer    User `json:"payer"`
	Receiver User `json:"receiver"`

	Type     int       `json:"type"`
	Ammount  float32   `json:"ammount"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
}

func (transaction *Transaction) SerializeJson() ([]byte, error) {
	data, err := json.MarshalIndent(*transaction, JsonPrefix, JsonIndent)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (transaction *Transaction) DeserializeJson(obj []byte) error {
	err := json.Unmarshal(obj, transaction)

	if err != nil {
		return err
	}

	return nil
}

func (transaction *Transaction) PopConstraints() {
	dbPayer, err := userservice.GetUser(transaction.Payer.Id)
	if err != nil {
		transaction.Payer.Expand(dbPayer)
	}

	dbReceiver, err := userservice.GetUser(transaction.Receiver.Id)
	if err != nil {
		transaction.Receiver.Expand(dbReceiver)
	}
}

func (transaction *Transaction) Expand(dbTransaction *dbmodels.Transaction) {
	transaction.Id = dbTransaction.Id
	transaction.Payer.Id = dbTransaction.PayerId
	transaction.Receiver.Id = dbTransaction.ReceiverId
	transaction.Type = dbTransaction.Type
	transaction.Ammount = dbTransaction.Ammount
	transaction.Currency = dbTransaction.Currency
	transaction.Date = dbTransaction.Date

	transaction.PopConstraints()
}

func (transaction *Transaction) Collapse() *dbmodels.Transaction {
	dbTransaction := dbmodels.Transaction{
		Id:         transaction.Id,
		PayerId:    transaction.Payer.Id,
		ReceiverId: transaction.Receiver.Id,
		Type:       transaction.Type,
		Ammount:    transaction.Ammount,
		Currency:   transaction.Currency,
		Date:       transaction.Date,
	}

	return &dbTransaction
}
