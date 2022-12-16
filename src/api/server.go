package api

import (
	"tok-core/src/definition"
)

func StartServer() {
	server := definition.Server

	// проставлять роуты
	setRoutes(server)

	// проставлять миддлвейры
	setMiddlewares(server)

	// запуск сервера
	server.Logger.Fatal(server.Start(definition.Config.Server.Port))
}