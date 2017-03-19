package httphandle

import (
	"gost/api"
	"gost/api/app/transactionapi"
	"gost/auth/identity"
	"gost/filter"
	"gost/orm/models"
	"gost/util/jsonutil"
	"net/http"

	"github.com/go-zoo/bone"
)

// Endpoint constants
const (
	transactionsGet    = "/transactions/{transactionId}"
	transactionsCreate = "/transactions"
)

// Authorization encompasses the identity provided by the auth package
type Authorization struct {
	Identity *identity.Identity
	Error    error
}

// InitializeRoutes initializes the application API routes and actions
func InitializeRoutes(mux *bone.Mux) {
	mux.HandleFunc(transactionsGet, func(rw http.ResponseWriter, req *http.Request) {
		requestAction(rw, req, http.MethodGet, transactionsGet, false, []string{identity.UserRoleNormal}, func(request *api.Request) api.Response {
			transactionID, found, err := filter.GetIDParameter("transactionId", request.Form)
			if err != nil {
				return api.BadRequest(err)
			}
			if !found {
				return api.NotFound(err)
			}

			return transactionapi.GetTransaction(transactionID)
		})
	})

	mux.HandleFunc(transactionsCreate, func(rw http.ResponseWriter, req *http.Request) {
		requestAction(rw, req, http.MethodGet, transactionsGet, false, []string{identity.UserRoleNormal}, func(request *api.Request) api.Response {
			transaction := &models.Transaction{}
			err := jsonutil.DeserializeJSON(request.Body, transaction)
			if err != nil {
				return api.BadRequest(api.ErrEntityFormat)
			}

			return transactionapi.CreateTransaction(transaction)
		})
	})
}
