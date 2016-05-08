package framework

import (
	"gost/config"
	"gost/util/jsonutil"
	"log"
)

const frameworkRoutes = `
[
    {
        "id": "AuthorizationRoute",
        "endpoint": "/auth",
        "actions": {
            "CreateSession": {
              "type": "POST",
              "allowAnonymous": true
            },
            "GetAllSessions": {
              "type": "GET",
              "allowAnonymous": false
            },
            "KillSession": {
              "type": "POST",
              "allowAnonymous": false
            },
            "ActivateAccount": {
              "type": "POST",
              "allowAnonymous": true
            },
            "ResendAccountActivationEmail": {
                "type": "GET",
                "allowAnonymous": true
            },
            "RequestResetPassword": {
              "type": "POST",
              "allowAnonymous": true
            },
            "ResetPassword": {
              "type": "POST",
              "allowAnonymous": true
            },
            "ChangePassword": {
              "type": "POST",
              "allowAnonymous": false
            }
        }
    },
    {
        "id": "ValuesRoute",
        "endpoint": "/values",
        "actions": {
            "Get": {
              "type": "GET",
              "allowAnonymous": false
            },
            "GetAnonymous": {
              "type": "GET",
              "allowAnonymous": true
            }
        }
    }
]`

// InitFrameworkRoutes initializes the routes used by the framework itself
func InitFrameworkRoutes() {
	var routes []config.Route

	var err = jsonutil.DeserializeJSON([]byte(frameworkRoutes), &routes)
	if err != nil {
		log.Fatalf("[InitFrameworkRoutes] %v\n", err)
	}

	err = config.AddRoutes(false, routes...)
	if err != nil {
		log.Fatalf("[InitFrameworkRoutes] %v\n", err)
	}
}
