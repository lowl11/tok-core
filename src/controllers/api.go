package controllers

import (
	"tok-core/src/controllers/auth_controller"
	"tok-core/src/controllers/profile_controller"
	"tok-core/src/controllers/static_controller"
	"tok-core/src/events"
	"tok-core/src/repositories"
)

type ApiControllers struct {
	Static  *static_controller.Controller
	Auth    *auth_controller.Controller
	Profile *profile_controller.Controller
}

func Get(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *ApiControllers {
	return &ApiControllers{
		Static: static_controller.Create(),
		Auth:   auth_controller.Create(apiRepositories, apiEvents),
	}
}
