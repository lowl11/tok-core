package controllers

import (
	"tok-core/src/controllers/auth_controller"
	"tok-core/src/controllers/feed_controller"
	"tok-core/src/controllers/notification_controller"
	"tok-core/src/controllers/post_controller"
	"tok-core/src/controllers/profile_controller"
	"tok-core/src/controllers/search_controller"
	"tok-core/src/controllers/static_controller"
	"tok-core/src/controllers/user_controller"
	"tok-core/src/events"
	"tok-core/src/repositories"
)

type ApiControllers struct {
	Static *static_controller.Controller

	// user flow
	Auth    *auth_controller.Controller
	Profile *profile_controller.Controller
	User    *user_controller.Controller

	// feed flow
	Feed *feed_controller.Controller
	Post *post_controller.Controller

	// search
	Search *search_controller.Controller

	// notification
	Notification *notification_controller.Controller
}

func Get(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *ApiControllers {
	return &ApiControllers{
		Static: static_controller.Create(),

		Auth:    auth_controller.Create(apiRepositories, apiEvents),
		Profile: profile_controller.Create(apiRepositories, apiEvents),
		User:    user_controller.Create(apiRepositories, apiEvents),

		Feed: feed_controller.Create(apiRepositories, apiEvents),
		Post: post_controller.Create(apiRepositories, apiEvents),

		Search: search_controller.Create(apiRepositories, apiEvents),

		Notification: notification_controller.Create(apiRepositories),
	}
}
