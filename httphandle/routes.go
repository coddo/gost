package httphandle

import (
	"gost/api"
	"gost/api/app/transactionapi"
	"gost/auth/identity"
	"gost/filter"
	"gost/orm/models"
	"gost/util/jsonutil"
	"net/http"
)

// Route represents an endpoint from the API
type Route struct {
	Path           string
	Method         string
	AllowAnonymous bool
	Roles          []string
	Action         func(request *api.Request) api.Response
}

// CreateAPIRoutes generates the main API routes used by the application
func CreateAPIRoutes() {
	Routes = append(Routes,
		&Route{
			Path:           "/transactions/{transactionId}",
			Method:         http.MethodGet,
			AllowAnonymous: false,
			Roles:          []string{identity.UserRoleNormal},
			Action:         getTransaction,
		},
		&Route{
			Path:           "/transactions",
			Method:         http.MethodPost,
			AllowAnonymous: false,
			Roles:          []string{identity.UserRoleNormal},
			Action:         createTransaction,
		},
	)
}

func getTransaction(request *api.Request) api.Response {
	transactionID, found, err := filter.GetIDParameter("transactionId", request.Form)
	if err != nil {
		return api.BadRequest(err)
	}
	if !found {
		return api.NotFound(err)
	}

	return transactionapi.GetTransaction(transactionID)
}

func createTransaction(request *api.Request) api.Response {
	transaction := &models.Transaction{}
	err := jsonutil.DeserializeJSON(request.Body, transaction)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	return transactionapi.CreateTransaction(transaction)
}
