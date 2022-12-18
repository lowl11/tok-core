package definition

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"tok-core/src/middlewares"
)

var Server *echo.Echo

func initServer() {
	Server = echo.New()

	Server.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	Server.Use(middleware.Secure())
	Server.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))
	Server.Use(middlewares.Timeout())
}
