package api

import (
	"tok-core/src/definition"
)

func StartServer() {
	server := definition.Server
	config := definition.Config.Server

	// проставлять роуты
	setRoutes(server)

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
