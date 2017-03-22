package devapi

import (
	"gost/api"
	"gost/auth"
	"net/http"
)

// AppUserModel is the model used for creating ApplicationUsers
type AppUserModel struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

// createAppUser is an endpoint used for creating application users
func createAppUser(model *AppUserModel) api.Response {
	user, err := auth.CreateAppUser(model.Email, model.Password, model.Roles)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.JSONResponse(http.StatusOK, user)
}

// activateAppUser is an endpoint for activating an app user
func activateAppUser(token string) api.Response {
	var err = auth.ActivateAppUser(token)
	if err != nil {
		return api.BadRequest(err)
	}

	return api.PlainTextResponse(http.StatusCreated, "Account is now active")
}
