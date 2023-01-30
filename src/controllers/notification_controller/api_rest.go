package notification_controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"tok-core/src/data/entities"
)

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
