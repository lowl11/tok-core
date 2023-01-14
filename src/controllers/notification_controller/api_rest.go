package notification_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/errors"
)

/*
	GetInfoREST обертка для _getInfo
*/
func (controller *Controller) GetInfoREST(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return controller.Error(ctx, errors.NotificationGetInfoParam)
	}

	info, err := controller._getInfo(username)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, info)
}

/*
	GetCountREST обертка для _getCount
*/
func (controller *Controller) GetCountREST(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return controller.Error(ctx, errors.NotificationGetCountParam)
	}

	count, err := controller._getCount(username)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, count)
}
