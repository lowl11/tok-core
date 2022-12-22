package api

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/middlewares"
)

// setMiddlewares общие миддлы вне зависимости от эндпоинта
func setMiddlewares(server *echo.Echo) {
	//
}

// setMiddlewaresPublic миддлы для публичных методов
func setMiddlewaresPublic(group *echo.Group) {
	group.Use(middlewares.Token)
	group.Use(middlewares.Session)
	group.Use(middlewares.IpAddress)
}

// setMiddlewaresPublicPart миддлы для публичных методов, но для запросов с частичными данными
func setMiddlewaresPublicPart(group *echo.Group) {
	group.Use(middlewares.TokenPart)
	group.Use(middlewares.SessionPart)
	group.Use(middlewares.IpAddress)
}
