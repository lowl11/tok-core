package user_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/definition"
)

func (controller *Controller) SignUp(ctx echo.Context) error {
	logger := definition.Logger

	logger.Info("SignUp")
	return controller.Ok(ctx, "OK")
}

func (controller *Controller) Login(ctx echo.Context) error {
	logger := definition.Logger

	logger.Info("Login")
	return controller.Ok(ctx, "OK")
}
