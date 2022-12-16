package controllers

import (
	"tok-core/src/controllers/static_controller"
	"tok-core/src/controllers/user_controller"
	"tok-core/src/events"
	"tok-core/src/repositories"
)

type ApiControllers struct {
	Static *static_controller.Controller
	User   *user_controller.Controller
}

func Get(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *ApiControllers {
	return &ApiControllers{
		Static: static_controller.Create(),
		User:   user_controller.Create(apiRepositories),
	}
}
