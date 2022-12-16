package profile_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) Update(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// привязка модели
	model := models.ProfileUpdate{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateBind.With(err))
	}

	// валидация модели
	if err := controller.validateProfileUpdate(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateValidate.With(err))
	}

	// изменение данных пользователя
	if err := controller.userRepo.UpdateProfile(session.Username, &model); err != nil {
		logger.Error(err, "Update profile info error")
		return controller.Error(ctx, errors.ProfileUpdate.With(err))
	}

	return controller.Ok(ctx, "OK")
}
