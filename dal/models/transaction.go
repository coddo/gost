package models

import (
	"gost/util/dateutil"
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

	PayerID    bson.ObjectId `json:"payer"`
	ReceiverID bson.ObjectId `json:"receiver"`

	PaymentPortal string `json:"paymentPortal"`
	PaymentToken  string `json:"paymentToken"`

	Type     int       `json:"type"`
	Ammount  float32   `json:"ammount"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
}

// Equal compares two Transaction objects. Implements the Objecter interface
func (transaction *Transaction) Equal(obj Objecter) bool {
	otherTransaction, ok := obj.(*Transaction)
	if !ok {
		return false
	}

	switch {
	case transaction.ID != otherTransaction.ID:
		return false
	case transaction.PayerID != otherTransaction.PayerID:
		return false
	case transaction.ReceiverID != otherTransaction.ReceiverID:
		return false
	case transaction.Type != otherTransaction.Type:
		return false
	case transaction.Ammount != otherTransaction.Ammount:
		return false
	case transaction.Currency != otherTransaction.Currency:
		return false
	case !dateutil.CompareDates(transaction.Date, otherTransaction.Date):
		return false
	}

	return true
}
