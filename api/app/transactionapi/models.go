package transactionapi

import (
	"gost/auth/identity"
	"gost/dal/models"
	"time"

	"gopkg.in/mgo.v2/bson"
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

// New creates a new transactions
func New() *Transaction {
	return new(Transaction)
}

// NewFrom creates a new transaction and copies the data from a DAL model
func NewFrom(dbTransaction *models.Transaction, payer, receiver *identity.ApplicationUser) *Transaction {
	return &Transaction{
		ID:            dbTransaction.ID,
		Payer:         payer,
		Receiver:      receiver,
		PaymentPortal: dbTransaction.PaymentPortal,
		PaymentToken:  dbTransaction.PaymentToken,
		Type:          dbTransaction.Type,
		Ammount:       dbTransaction.Ammount,
		Currency:      dbTransaction.Currency,
		Date:          dbTransaction.Date,
	}
}

// ToDalModel converts a transaction to a DAL model that contains the respective data
func (transaction *Transaction) ToDalModel() *models.Transaction {
	return &models.Transaction{
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
}
