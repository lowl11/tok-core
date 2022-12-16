package profile_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/repositories"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base
	userRepo *user_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories) *Controller {
	return &Controller{
		userRepo: apiRepositories.User,
	}
}
