package apifilter

import (
	"go-server-template/models"
)

func CheckTransactionIntegrity(transaction *models.Transaction) bool {
	switch {
	case len(transaction.Payer.Id) == 0:
		return false
	case len(transaction.Receiver.Id) == 0:
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
