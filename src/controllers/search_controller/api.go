package search_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
	"tok-core/src/services/query_corrector"
	"tok-core/src/services/search_sorter"
)

/*
	Smart поиск по нескольким типам данных
	Комбинируются несколько типов данных и отдаются клиенту в перемешке
*/
func (controller *Controller) Smart(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// связка модели
	model := models.SearchSmart{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SearchSmartBind.With(err))
	}

	// валидация модели
	if err := controller.validateSmart(&model); err != nil {
		return controller.Error(ctx, errors.SearchSmartValidate.With(err))
	}

	// обработка поискового запроса клиента
	model.Query = query_corrector.Correct(model.Query)

	// получаем пользователей
	// поиск пользователей по запросу клиента
	users, err := controller.userRepo.Search(model.Query)
	if err != nil {
		logger.Error(err, "Search users error")
		return controller.Error(ctx, errors.SearchUser.With(err))
	}

	// определяем подписан ли клиент на список пользователей
	mySubscriptions := array.NewWithList[string](session.Subscriptions.Subscriptions...)

	// обработка списка пользователей постов для клиента
	userList := make([]models.SearchUserGet, 0, len(users))
	for _, user := range users {
		userList = append(userList, models.SearchUserGet{
			Username:   user.Username,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Subscribed: mySubscriptions.Contains(user.Username),
		})
	}

	// получить список категорий постов
	categories, err := controller.postCategoryRepo.Search(model.Query)
	if err != nil {
		logger.Error(err, "Search categories error")
		return controller.Error(ctx, errors.SearchCategory.With(err))
	}

	// обработка списка категорий постов для клиента
	categoryList := make([]models.SearchCategoryGet, 0, len(categories))
	for _, item := range categories {
		categoryList = append(categoryList, models.SearchCategoryGet{
			Code: item.Code,
			Name: item.Name,
		})
	}

	// сортировка и компановка двух листов в один
	list := search_sorter.SortSmart(userList, categoryList)

	return controller.Ok(ctx, list)
}

/*
	User поиск пользователей
*/
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
	list := make([]models.SearchUserGet, 0, len(users))
	for _, user := range users {
		list = append(list, models.SearchUserGet{
			Username:   user.Username,
			Name:       user.Name,
			Avatar:     user.Avatar,
			Subscribed: mySubscriptions.Contains(user.Username),
		})
	}

	return controller.Ok(ctx, list)
}

/*
	Category поиск по категориям
*/
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

	if model.Query == "" {
		categories, err := controller.postCategoryRepo.GetFirstSorted()
		if err != nil {
			return controller.Error(ctx, errors.PostCategoryGetList.With(err))
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
