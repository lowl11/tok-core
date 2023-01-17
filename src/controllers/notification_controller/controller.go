package notification_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/repositories"
	"tok-core/src/repositories/notification_repository"
	"tok-core/src/repositories/post_repository"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base

	userRepo         *user_repository.Repository
	postRepo         *post_repository.Repository
	notificationRepo *notification_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories) *Controller {
	return &Controller{
		userRepo:         apiRepositories.User,
		postRepo:         apiRepositories.Post,
		notificationRepo: apiRepositories.Notification,
	}
}
