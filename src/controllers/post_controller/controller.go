package post_controller

import (
	"tok-core/src/controllers/controller"
	"tok-core/src/events"
	"tok-core/src/events/image_event"
	"tok-core/src/repositories"
	"tok-core/src/repositories/category_count_repository"
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
	userInterest      *user_interest_repository.Repository

	image *image_event.Event
}

func Create(apiRepositories *repositories.ApiRepositories, apiEvents *events.ApiEvents) *Controller {
	return &Controller{
		userRepo: apiRepositories.User,

		postRepo:         apiRepositories.Post,
		postCategoryRepo: apiRepositories.PostCategory,
		postCommentRepo:  apiRepositories.PostComment,
		postLikeRepo:     apiRepositories.PostLike,

		categoryCountRepo: apiRepositories.CategoryCount,
		userInterest:      apiRepositories.UserInterest,

		image: apiEvents.Image,
	}
}
