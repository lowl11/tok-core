package search_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
	"tok-core/src/services/query_corrector"
)

func (controller *Controller) Smart(ctx echo.Context) error {
	return controller.Ok(ctx, "OK")
}

func (controller *Controller) User(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// связка модели
	model := models.SearchUser{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SearchUserBind.With(err))
	}

	// валидация модели
	if err := controller.validateUser(&model); err != nil {
		return controller.Error(ctx, errors.SearchUserValidate.With(err))
	}

	// обработка поискового запроса клиента
	model.Query = query_corrector.Correct(model.Query)

	// поиск пользователей по запросу клиента
	users, err := controller.userRepo.Search(model.Query)
	if err != nil {
		logger.Error(err, "Search users error")
		return controller.Error(ctx, errors.SearchUser.With(err))
	}

	// определяем подписан ли клиент на список пользователей
	mySubscriptions := array.NewWithList[string](session.Subscriptions.Subscriptions...)

	// обработка списка пользователей постов для клиента
	list := make([]models.UserSearch, 0, len(users))
	for _, user := range users {
		list = append(list, models.UserSearch{
			Username:   user.Username,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Subscribed: mySubscriptions.Contains(user.Username),
		})
	}

	return controller.Ok(ctx, list)
}

func (controller *Controller) Category(ctx echo.Context) error {
	logger := definition.Logger

	// связка модели
	model := models.SearchCategory{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SearchCategoryBind.With(err))
	}

	// валидация модели
	if err := controller.validateCategory(&model); err != nil {
		return controller.Error(ctx, errors.SearchCategoryValidate.With(err))
	}

	// обработка поискового запроса клиента
	model.Query = query_corrector.Correct(model.Query)

	// получить список категорий постов
	categories, err := controller.postCategoryRepo.Search(model.Query)
	if err != nil {
		logger.Error(err, "Search categories error")
		return controller.Error(ctx, errors.SearchCategory.With(err))
	}

	// обработка списка категорий постов для клиента
	list := make([]models.PostCategoryGet, 0, len(categories))
	for _, item := range categories {
		list = append(list, models.PostCategoryGet{
			Code: item.Code,
			Name: item.Name,
		})
	}

	return controller.Ok(ctx, list)
}
