package httphandletest

import (
	"gost/httphandle"
	"gost/servers"
)

// InitTestsRoutes initializez the routes used for testing the endpoints
func InitTestsRoutes() {
	httphandle.CreateFrameworkRoutes()
	httphandle.CreateDevelopmentRoutes()
	httphandle.CreateAPIRoutes()
	httphandle.InitRoutes(servers.Multiplexer)
}
