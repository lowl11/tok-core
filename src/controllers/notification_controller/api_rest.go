package notification_controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
)

/*
	ReadREST обертка для _getCount
*/
func (controller *Controller) ReadREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.NotificationRead{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.NotificationReadBind.With(err))
	}

	if err := controller.validationRead(&model); err != nil {
		return controller.Error(ctx, errors.NotificationReadValidation.With(err))
	}

	if err := controller._read(session.Username); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	GetInfoREST обертка для _getInfo
*/
func (controller *Controller) GetInfoREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	info, err := controller._getInfo(session.Username, page)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, info)
}

/*
	GetCountREST обертка для _getCount
*/
func (controller *Controller) GetCountREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	count, err := controller._getCount(session.Username)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, count)
}
