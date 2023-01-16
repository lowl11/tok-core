package post_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/events/image_event"
	"tok-core/src/events/notification_event"
	"tok-core/src/events/post_category_event"
	"tok-core/src/repositories"
	"tok-core/src/repositories/category_count_repository"
	"tok-core/src/repositories/feed_repository"
	"tok-core/src/repositories/post_category_repository"
	"tok-core/src/repositories/post_comment_repository"
	"tok-core/src/repositories/post_like_repository"
	"tok-core/src/repositories/post_repository"
	"tok-core/src/repositories/user_interest_repository"
	"tok-core/src/repositories/user_repository"
)

type Controller struct {
	controller.Base

	userRepo *user_repository.Repository

	postRepo         *post_repository.Repository
	postCategoryRepo *post_category_repository.Repository
	postCommentRepo  *post_comment_repository.Repository
	postLikeRepo     *post_like_repository.Repository

	categoryCountRepo *category_count_repository.Repository
	userInterestRepo  *user_interest_repository.Repository
	feedRepo          *feed_repository.Repository

	image        *image_event.Event
	category     *post_category_event.Event
	notification *notification_event.Event
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		userRepo: apiRepositories.User,

		postRepo:         apiRepositories.Post,
		postCategoryRepo: apiRepositories.PostCategory,
		postCommentRepo:  apiRepositories.PostComment,
		postLikeRepo:     apiRepositories.PostLike,

		categoryCountRepo: apiRepositories.CategoryCount,
		userInterestRepo:  apiRepositories.UserInterest,
		feedRepo:          apiRepositories.Feed,

		image:        apiEvents.Image,
		category:     apiEvents.PostCategory,
		notification: apiEvents.Notification,
	}
}
