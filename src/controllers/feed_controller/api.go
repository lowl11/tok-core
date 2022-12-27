package feed_controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

/*
	Main лента на главной странице
*/
func (controller *Controller) Main(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// чтение номера страницы
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1 // номер страницы по умолчанию
	}

	// список подписок из сессии
	subscriptions := session.Subscriptions.Subscriptions

	// добавить себя в список
	subscriptions = append(subscriptions, session.Username)

	// посты с массивом юзернеймов из подписок
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByUsernameList(subscriptions, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by username list error")
		return controller.Error(ctx, errors.PostsGetByUsernameList.With(err))
	}

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

/*
	Explore лента "рекомендаций"
*/
func (controller *Controller) Explore(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// чтение номера страницы
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1 // номер страницы по умолчанию
	}

	// список подписок из сессии
	subscriptions := session.Subscriptions.Subscriptions

	// добавить себя в список
	subscriptions = append(subscriptions, session.Username)

	// посты с массивом юзернеймов из подписок
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByUsernameList(subscriptions, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by username list error")
		return controller.Error(ctx, errors.PostsGetByUsernameList.With(err))
	}

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

/*
	User лента на странице пользователя
*/
func (controller *Controller) User(ctx echo.Context) error {
	logger := definition.Logger

	// чтение параметра логина пользователя
	username := ctx.Param("username")
	if username == "" {
		return controller.Error(ctx, errors.PostsGetByUsernameParam)
	}

	// чтение номера страницы
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1 // номер страницы по умолчанию
	}

	// получить список постовов по логину
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByUsername(username, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by username error")
		return controller.Error(ctx, errors.PostsGetByUsername.With(err))
	}

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

/*
	Category лента на странице категории
*/
func (controller *Controller) Category(ctx echo.Context) error {
	logger := definition.Logger

	// чтение номера страницы
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1 // номер страницы по умолчанию
	}

	// чтение параметра кода категории
	categoryCode := ctx.Param("category_code")
	if categoryCode == "" {
		return controller.Error(ctx, errors.PostsGetByCategoryParam)
	}

	// получить список постовов по логину
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByCategory(categoryCode, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by category error")
		return controller.Error(ctx, errors.PostsGetByCategory.With(err))
	}

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,
		})
	}

	return controller.Ok(ctx, list)
}

/*
	Single страница одного поста
*/
func (controller *Controller) Single(ctx echo.Context) error {
	logger := definition.Logger

	// чтение параметра кода поста
	code := ctx.Param("code")
	if code == "" {
		return controller.Error(ctx, errors.PostsGetBySingleParam)
	}

	// получение списка постов
	post, err := controller.postRepo.GetByCode(code)
	if err != nil {
		logger.Error(err, "Get post by code error")
		return controller.Error(ctx, errors.PostsGetByCode.With(err))
	}

	// если пост не найден
	if post == nil {
		return controller.NotFound(ctx, errors.PostNotFound)
	}

	// обработка поста для клиента
	return controller.Ok(ctx, &models.PostGet{
		AuthorUsername: post.AuthorUsername,
		AuthorName:     post.AuthorName,
		AuthorAvatar:   post.AuthorAvatar,

		CategoryCode: post.CategoryCode,
		CategoryName: post.CategoryName,

		Code: post.Code,
		Text: post.Text,
		Picture: &models.PostGetPicture{
			Path:   post.Picture,
			Width:  post.PictureWidth,
			Height: post.PictureHeight,
		},
		CreatedAt: post.CreatedAt,
	})
}
