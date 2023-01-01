package feed_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/events/feed_event"
	"tok-core/src/repositories"
	"tok-core/src/repositories/post_like_repository"
	"tok-core/src/repositories/post_repository"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base

	postRepo     *post_repository.Repository
	postLikeRepo *post_like_repository.Repository
	userRepo     *user_repository.Repository

	feed *feed_event.Event
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		postRepo:     apiRepositories.Post,
		postLikeRepo: apiRepositories.PostLike,
		userRepo:     apiRepositories.User,

		feed: apiEvents.Feed,
	}
}
