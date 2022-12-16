package api

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/controllers"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/repositories"
)

func setRoutes(server *echo.Echo) {
	logger := definition.Logger

	// ивенты
	apiEvents, err := events.Get()
	if err != nil {
		logger.Fatal(err, "Initializing events error")
	}

	// репозитории
	apiRepositories, err := repositories.Get(apiEvents)
	if err != nil {
		logger.Fatal(err, "Connecting to database error")
	}

	// контроллеры
	apiControllers := controllers.Get(apiRepositories, apiEvents)

	// статичные методы
	server.GET("/health", apiControllers.Static.Health)
	server.RouteNotFound("*", apiControllers.Static.RouteNotFound)
}