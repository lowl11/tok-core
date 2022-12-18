package post_controller

import "github.com/labstack/echo/v4"

func (controller *Controller) Add(ctx echo.Context) error {
	return controller.Ok(ctx, "OK")
}
