package profile_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/events/client_session_event"
	"tok-core/src/events/image_event"
	"tok-core/src/repositories"
	"tok-core/src/repositories/subscription_repository"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base

	userRepo      *user_repository.Repository
	clientSession *client_session_event.Event

	image         *image_event.Event
	subscriptRepo *subscription_repository.Repository
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		userRepo:      apiRepositories.User,
		subscriptRepo: apiRepositories.Subscription,

		clientSession: apiEvents.ClientSession,
		image:         apiEvents.Image,
	}
}
