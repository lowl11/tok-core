package feed_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"github.com/lowl11/lazy-collection/set"
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazylog/layers"
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

	likeList, err := controller.postLikeRepo.GetByList(type_list.NewWithList[entities.PostGet, string](posts...).Select(func(item entities.PostGet) string {
		return item.Code
	}).Slice())
	if err != nil {
		logger.Error(err, "Get posts like counter error", layers.Mongo)
	}

	var likeArray *array.Array[entities.PostLikeGetList]
	if likeList != nil {
		likeArray = array.NewWithList[entities.PostLikeGetList](likeList...)
	}

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		var myLike bool
		var likeCount int

		if likeArray != nil {
			foundLike := likeArray.Single(func(iterator entities.PostLikeGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundLike != nil {
				likeCount = foundLike.LikesCount

				likeAuthors := array.NewWithList[string](foundLike.LikeAuthors...)
				myLike = likeAuthors.Contains(session.Username)
			}
		}

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

			LikeCount: likeCount,
			MyLike:    myLike,
		})
	}

	return controller.Ok(ctx, list)
}

/*
	_explore лента "рекомендаций"
*/
func (controller *Controller) _explore(session *entities.ClientSession, page int) ([]models.PostGet, *models.Error) {
	logger := definition.Logger

	// запрос для получения "рекомендаций"
	posts, err := controller.feed.GetExplore(session.Username, []string{"Мемы", "Образование", "memy", "obrazovanie"}, page)
	if err != nil {
		logger.Error(err, "Get list for explore error", layers.Elastic)
		return nil, errors.PostsGetExplore.With(err)
	}

	usernames := set.NewWithSize[string](len(posts))
	for _, post := range posts {
		usernames.Push(post.Author)
	}

	// получаем "имена" и "аватары" авторов
	dynamics, err := controller.userRepo.GetDynamicByUsernames(usernames.Slice())
	if err != nil {
		return nil, errors.UserDynamicGet.With(err)
	}

	dynamicsArray := array.NewWithList[entities.UserDynamicGet](dynamics...)

	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		var authorName *string
		var authorAvatar *string

		dynamicUserInfo := dynamicsArray.Single(func(item entities.UserDynamicGet) bool {
			return item.Username == item.Username
		})
		if dynamicUserInfo != nil {
			authorName = dynamicUserInfo.Name
			authorAvatar = dynamicUserInfo.Avatar
		}

		list = append(list, models.PostGet{
			AuthorUsername: item.Author,
			AuthorName:     authorName,
			AuthorAvatar:   authorAvatar,

			CategoryCode: item.Category,
			CategoryName: item.CategoryName,

			Code:      item.Code,
			Text:      item.Text,
			Picture:   item.Picture,
			CreatedAt: item.CreatedAt,
		})
	}

	return list, nil
}

/*
	ExploreREST обертка для _explore
*/
func (controller *Controller) ExploreREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	// чтение номера страницы
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1 // номер страницы по умолчанию
	}

	list, err := controller._explore(session, page)
	if err != nil {
		return controller.Error(ctx, err)
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
