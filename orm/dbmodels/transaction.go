package dbmodels

import (
	"gost/util"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Transaction is a struct representing transactions between users
type Transaction struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`

	PayerID    bson.ObjectId `bson:"payerId,omitempty" json:"payerId"`
	ReceiverID bson.ObjectId `bson:"receiverId,omitempty" json:"receiverId"`

	PaymentPortal string `bson:"paymentPortal,omitempty" json:"paymentPortal"`
	PaymentToken  string `bson:"paymentToken,omitempty" json:"paymentToken"`

	Type     int       `bson:"type,omitempty" json:"type"`
	Ammount  float32   `bson:"ammount,omitempty" json:"ammount"`
	Currency string    `bson:"currency,omitempty" json:"currency"`
	Date     time.Time `bson:"date,omitempty" json:"date"`
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
	case !util.CompareDates(transaction.Date, otherTransaction.Date):
		return false
	}

	return true
}
