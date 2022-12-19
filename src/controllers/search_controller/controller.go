package search_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/repositories"
	"tok-core/src/repositories/post_category_repository"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base

	userRepo         *user_repository.Repository
	postCategoryRepo *post_category_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories, _ *events.ApiEvents) *Controller {
	return &Controller{
		userRepo:         apiRepositories.User,
		postCategoryRepo: apiRepositories.PostCategory,
	}
}
