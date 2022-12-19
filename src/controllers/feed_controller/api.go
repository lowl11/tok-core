package feed_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) Main(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// список подписок из сессии
	subscriptions := session.Subscriptions.Subscriptions

	// посты с массивом юзернеймов из подписок
	posts, err := controller.postRepo.GetByUsernameList(subscriptions)
	if err != nil {
		logger.Error(err, "Get posts list by username list error")
		return controller.Error(ctx, errors.PostsGetByUsernameList.With(err))
	}

	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code:      item.Code,
			Text:      item.Text,
			Picture:   item.Picture,
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

func (controller *Controller) User(ctx echo.Context) error {
	logger := definition.Logger

	// query
	username := ctx.Param("username")
	if username == "" {
		return controller.Error(ctx, errors.PostsGetByUsernameParam)
	}

	// получить список постовов по логину
	posts, err := controller.postRepo.GetByUsername(username)
	if err != nil {
		logger.Error(err, "Get posts list by username error")
		return controller.Error(ctx, errors.PostsGetByUsername.With(err))
	}

	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code:      item.Code,
			Text:      item.Text,
			Picture:   item.Picture,
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

func (controller *Controller) Category(ctx echo.Context) error {
	logger := definition.Logger

	// query
	categoryCode := ctx.Param("category_code")
	if categoryCode == "" {
		return controller.Error(ctx, errors.PostsGetByCategoryParam)
	}

	// получить список постовов по логину
	posts, err := controller.postRepo.GetByCategory(categoryCode)
	if err != nil {
		logger.Error(err, "Get posts list by category error")
		return controller.Error(ctx, errors.PostsGetByCategory.With(err))
	}

	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code:      item.Code,
			Text:      item.Text,
			Picture:   item.Picture,
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

func (controller *Controller) Single(ctx echo.Context) error {
	logger := definition.Logger

	code := ctx.Param("code")
	if code == "" {
		return controller.Error(ctx, errors.PostsGetBySingleParam)
	}

	post, err := controller.postRepo.GetByCode(code)
	if err != nil {
		logger.Error(err, "Get post by code error")
		return controller.Error(ctx, errors.PostsGetByCode.With(err))
	}

	if post == nil {
		return controller.Error(ctx, errors.PostNotFound)
	}

	return controller.Ok(ctx, &models.PostGet{
		AuthorUsername: post.AuthorUsername,
		AuthorName:     post.AuthorName,
		AuthorAvatar:   post.AuthorAvatar,

		CategoryCode: post.CategoryCode,
		CategoryName: post.CategoryName,

		Code:      post.Code,
		Text:      post.Text,
		Picture:   post.Picture,
		CreatedAt: post.CreatedAt,
	})
}
