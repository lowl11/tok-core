package feed_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/repositories"
	"tok-core/src/repositories/post_repository"
)

type Controller struct {
	controller.Base

	postRepo *post_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		postRepo: apiRepositories.Post,
	}
}
