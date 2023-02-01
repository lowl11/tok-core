package definition

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"tok-core/src/middlewares"
)

var Server *echo.Echo

// initServer создание объекта сервера
func initServer() {
	Server = echo.New()

	// общие миддлы
	Server.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	Server.Use(middleware.Secure())
	Server.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))
	Server.Use(middlewares.Timeout())
	Server.Use(middlewares.IpAddress)
}
