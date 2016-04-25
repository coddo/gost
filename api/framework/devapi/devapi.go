package devapi

import (
	"gost/api"
	"gost/auth"
	"gost/config"
	"gost/filter"
	"gost/util"
	"net/http"
)

// DevAPI defines the API endpoint for development actions and custom testing
type DevAPI int

// AppUserModel is the model used for creating ApplicationUsers
type AppUserModel struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccountType int    `json:"accountType"` // 0 - NormalUser | 1 - Admin
}

// CreateAppUser is an endpoint used for creating application users
func (v *DevAPI) CreateAppUser(params *api.Request) api.Response {
	model := &AppUserModel{}

	err := util.DeserializeJSON(params.Body, model)
	if err != nil {
		return api.BadRequest(api.ErrEntityFormat)
	}

	var activationServiceLink = config.HTTPServerAddress + config.APIInstance + "dev/ActivateAppUser?token=%s"

	user, err := auth.CreateAppUser(model.Email, model.Password, model.AccountType, activationServiceLink)
	if err != nil {
		return api.InternalServerError(err)
	}

	return api.JSONResponse(http.StatusOK, user)
}

// ActivateAppUser is an endpoint for activating an app user
func (v *DevAPI) ActivateAppUser(params *api.Request) api.Response {
	var token, found = filter.GetStringParameter("token", params.Form)
	if !found {
		return api.BadRequest(api.ErrInvalidInput)
	}

	var err = auth.ActivateAppUser(token)
	if err != nil {
		return api.BadRequest(err)
	}

	return api.PlainTextResponse(http.StatusOK, "Account is now active")
}
