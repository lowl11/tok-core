package auth_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/events/client_session_event"
	"tok-core/src/repositories"
	"tok-core/src/repositories/auth_repository"
	"tok-core/src/repositories/subscription_repository"
	"tok-core/src/repositories/user_ip_repository"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base

	authRepo      *auth_repository.Repository
	userRepo      *user_repository.Repository
	userIpRepo    *user_ip_repository.Repository
	subscriptRepo *subscription_repository.Repository

	clientSession *client_session_event.Event
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		authRepo:      apiRepositories.Auth,
		userRepo:      apiRepositories.User,
		subscriptRepo: apiRepositories.Subscription,

		clientSession: apiEvents.ClientSession,
	}
}
