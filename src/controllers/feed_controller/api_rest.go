package feed_controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
)

/*
	MainREST обертка для _main
*/
func (controller *Controller) MainREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	// чтение номера страницы
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1 // номер страницы по умолчанию
	}

	list, err := controller._main(session, page)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, list)
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
	UserREST обертка для _user
*/
func (controller *Controller) UserREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

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

	list, err := controller._user(session, username, page)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, list)
}

/*
	CategoryREST обертка для _category
*/
func (controller *Controller) CategoryREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

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

	list, err := controller._category(session, categoryCode, page)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, list)
}

/*
	SingleREST обертка для _single
*/
func (controller *Controller) SingleREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	// чтение параметра кода поста
	code := ctx.Param("code")
	if code == "" {
		return controller.Error(ctx, errors.PostsGetBySingleParam)
	}

	post, err := controller._single(session, code)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, post)
}
