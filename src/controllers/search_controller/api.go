package search_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
)

func (controller *Controller) Smart(ctx echo.Context) error {
	return controller.Ok(ctx, "OK")
}

func (controller *Controller) User(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	// query
	model := models.SearchUser{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SearchUserBind.With(err))
	}

	if err := controller.validateUser(&model); err != nil {
		return controller.Error(ctx, errors.SearchUserValidate.With(err))
	}

	// TODO: обработать query (trim, lowercase, etc.)

	users, err := controller.userRepo.Search(model.Query)
	if err != nil {
		return controller.Error(ctx, errors.UserSearch.With(err))
	}

	mySubscriptions := array.NewWithList[string](session.Subscriptions.Subscriptions...)

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
	return controller.Ok(ctx, "OK")
}
