package main

import (
	"tok-core/src/api"
	"tok-core/src/controllers"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/job"
	"tok-core/src/repositories"
)

func main() {
	definition.Init()
	logger := definition.Logger

	// инициализация ивентов
	apiEvents, err := events.Get()
	if err != nil {
		logger.Fatal(err, "Initializing events error", "Application")
	}

	// инициализация репозиториев
	apiRepositories, err := repositories.Get(apiEvents)
	if err != nil {
		logger.Fatal(err, "Initializing repositories error", "Application")
	}

	// контроллеры
	apiControllers := controllers.Get(apiRepositories, apiEvents)

	// запуск джобов
	job.RunAsync(apiControllers)

	// запуск сервера
	api.StartServer(apiControllers, apiRepositories, apiEvents)
}
