package api

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/middlewares"
)

func setMiddlewares(server *echo.Echo) {
	//
}

func setMiddlewaresToGroup(group *echo.Group) {
	//
}

func setMiddlewaresPublic(group *echo.Group) {
	group.Use(middlewares.Session)
}
