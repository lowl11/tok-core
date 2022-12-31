package post_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/events/feed_event"
	"tok-core/src/events/image_event"
	"tok-core/src/repositories"
	"tok-core/src/repositories/post_category_repository"
	"tok-core/src/repositories/post_comment_repository"
	"tok-core/src/repositories/post_like_repository"
	"tok-core/src/repositories/post_repository"
)

type Controller struct {
	controller.Base

	postRepo         *post_repository.Repository
	postCategoryRepo *post_category_repository.Repository
	postCommentRepo  *post_comment_repository.Repository
	postLikeRepo     *post_like_repository.Repository

	feed  *feed_event.Event
	image *image_event.Event
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		postRepo:         apiRepositories.Post,
		postCategoryRepo: apiRepositories.PostCategory,
		postCommentRepo:  apiRepositories.PostComment,
		postLikeRepo:     apiRepositories.PostLike,

		feed:  apiEvents.Feed,
		image: apiEvents.Image,
	}
}
