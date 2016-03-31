package models

import (
	"gost/auth/identity"
	"gost/orm/dbmodels"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Constants representing the type of transaction
const (
	TransactionTypeCash = iota
	TransactionTypeCard = iota
)

// Transaction is a struct representing transactions between users
type Transaction struct {
	ID bson.ObjectId `json:"id"`

	Payer    *identity.ApplicationUser `json:"payer"`
	Receiver *identity.ApplicationUser `json:"receiver"`

	PaymentPortal string `json:"paymentPortal"`
	PaymentToken  string `json:"paymentToken"`

	Type     int       `json:"type"`
	Ammount  float32   `json:"ammount"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
}

// Expand copies the dbmodels.Transaction to a Transaction expands all
// the components by fetching them from the database
func (transaction *Transaction) Expand(dbTransaction *dbmodels.Transaction) {
	transaction.ID = dbTransaction.ID
	transaction.Payer.ID = dbTransaction.PayerID
	transaction.Receiver.ID = dbTransaction.ReceiverID
	transaction.PaymentPortal = dbTransaction.PaymentPortal
	transaction.PaymentToken = dbTransaction.PaymentToken
	transaction.Type = dbTransaction.Type
	transaction.Ammount = dbTransaction.Ammount
	transaction.Currency = dbTransaction.Currency
	transaction.Date = dbTransaction.Date
}

// Collapse coppies the Transaction to a dbmodels.Transaction user and
// only keeps the unique identifiers from the inner components
func (transaction *Transaction) Collapse() *dbmodels.Transaction {
	dbTransaction := dbmodels.Transaction{
		ID:            transaction.ID,
		PayerID:       transaction.Payer.ID,
		ReceiverID:    transaction.Receiver.ID,
		PaymentPortal: transaction.PaymentPortal,
		PaymentToken:  transaction.PaymentToken,
		Type:          transaction.Type,
		Ammount:       transaction.Ammount,
		Currency:      transaction.Currency,
		Date:          transaction.Date,
	}

	return &dbTransaction
}
