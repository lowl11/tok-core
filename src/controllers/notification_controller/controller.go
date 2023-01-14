package notification_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/repositories"
	"tok-core/src/repositories/notification_repository"
)

type Controller struct {
	controller.Base

	notificationRepo *notification_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories) *Controller {
	return &Controller{
		notificationRepo: apiRepositories.Notification,
	}
}
