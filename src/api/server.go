package api

import (
	"tok-core/src/controllers"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/repositories"
)

func StartServer(apiControllers *controllers.ApiControllers, apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) {
	server := definition.Server
	config := definition.Config.Server

	// проставлять роуты
	setRoutes(server, apiControllers, apiRepositories, apiEvents)

	// проставлять миддлвейры
	setMiddlewares(server)

	// запуск сервера
	var port string
	if definition.Config.Primary {
		port = config.Port.Primary
	} else {
		port = config.Port.Secondary
	}
	server.Logger.Fatal(server.Start(port))
}
