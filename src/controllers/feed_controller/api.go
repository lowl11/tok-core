package feed_controller

import "github.com/labstack/echo/v4"

func (controller *Controller) Main(ctx echo.Context) error {
	return controller.Ok(ctx, "OK")
}

func (controller *Controller) User(ctx echo.Context) error {
	return controller.Ok(ctx, "OK")
}

func (controller *Controller) Category(ctx echo.Context) error {
	return controller.Ok(ctx, "OK")
}
