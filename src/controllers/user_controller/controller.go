package user_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/repositories"
	"tok-core/src/repositories/auth_repository"
)

type Controller struct {
	controller.Base
	authRepo *auth_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories) *Controller {
	return &Controller{
		authRepo: apiRepositories.Auth,
	}
}
