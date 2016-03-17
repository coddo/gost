package models

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"gost/service/appuserservice"
	"time"
)

type Transaction struct {
	Id bson.ObjectId `json:"id"`

	Payer    ApplicationUser `json:"payer"`
	Receiver ApplicationUser `json:"receiver"`

	PaymentPortal string `json:"paymentPortal"`
	PaymentToken  string `json:"paymentToken"`

	Type     int       `json:"type"`
	Ammount  float32   `json:"ammount"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
}

func (transaction *Transaction) PopConstrains() {
	dbPayer, err := appuserservice.GetUser(transaction.Payer.Id)
	if err != nil {
		transaction.Payer.Expand(dbPayer)
	}

	dbReceiver, err := appuserservice.GetUser(transaction.Receiver.Id)
	if err != nil {
		transaction.Receiver.Expand(dbReceiver)
	}
}

func (transaction *Transaction) Expand(dbTransaction *dbmodels.Transaction) {
	transaction.Id = dbTransaction.ID
	transaction.Payer.Id = dbTransaction.PayerID
	transaction.Receiver.Id = dbTransaction.ReceiverID
	transaction.PaymentPortal = dbTransaction.PaymentPortal
	transaction.PaymentToken = dbTransaction.PaymentToken
	transaction.Type = dbTransaction.Type
	transaction.Ammount = dbTransaction.Ammount
	transaction.Currency = dbTransaction.Currency
	transaction.Date = dbTransaction.Date

	transaction.PopConstrains()
}

func (transaction *Transaction) Collapse() *dbmodels.Transaction {
	dbTransaction := dbmodels.Transaction{
		ID:            transaction.Id,
		PayerID:       transaction.Payer.Id,
		ReceiverID:    transaction.Receiver.Id,
		PaymentPortal: transaction.PaymentPortal,
		PaymentToken:  transaction.PaymentToken,
		Type:          transaction.Type,
		Ammount:       transaction.Ammount,
		Currency:      transaction.Currency,
		Date:          transaction.Date,
	}

	return &dbTransaction
}
