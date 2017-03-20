package devapi

import (
	"fmt"
	"gost/api"
	"gost/auth"
	"gost/config"
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
	var activationServiceLink = fmt.Sprintf("%s://%s%s%s", config.ServerType, config.HTTPServerAddress, config.APIInstance, "dev/ActivateAppUser?token=%s")

	user, err := auth.CreateAppUser(model.Email, model.Password, model.Roles, activationServiceLink)
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

	return api.PlainTextResponse(http.StatusOK, "Account is now active")
}
