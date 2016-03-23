package apifilter

import (
	"gost/orm/models"
)

// CheckTransactionIntegrity checks if a Transaction has all the compulsory fields populated
func CheckTransactionIntegrity(transaction *models.Transaction) bool {
	switch {
	case len(transaction.Payer.ID) == 0:
		return false
	case len(transaction.Receiver.ID) == 0:
		return false
	case transaction.Type < 0:
		return false
	case transaction.Ammount < 0:
		return false
	case len(transaction.Currency) == 0:
		return false
	}

	return true
}
