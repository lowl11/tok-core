package profile_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
)

/*
	Subscribe REST обертка для _subscribe
*/
func (controller *Controller) SubscribeREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	// привязка модели
	model := models.ProfileSubscribe{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileBind.With(err))
	}

	// валидация модели
	if err := controller.validateSubscribeProfile(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileValidate.With(err))
	}

	if err := controller._subscribe(session, token, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	Unsubscribe REST обертка для _unsubscribe
*/
func (controller *Controller) UnsubscribeREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	// привязка модели
	model := models.ProfileUnsubscribe{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileBind.With(err))
	}

	// валидация модели
	if err := controller.validateUnsubscribeProfile(&model); err != nil {
		return controller.Error(ctx, errors.UnsubscribeOfProfileValidate.With(err))
	}

	// отписка
	if err := controller._unsubscribe(session, token, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	Update REST обертка для _update
*/
func (controller *Controller) UpdateREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	// привязка модели
	model := models.ProfileUpdate{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateBind.With(err))
	}

	// валидация модели
	if err := controller.validateUpdate(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateValidate.With(err))
	}

	if err := controller._update(session, token, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	UpdateContacts REST обертка для _updateContacts
*/
func (controller *Controller) UpdateContactsREST(ctx echo.Context) error {
	// привязка модели
	model := models.ProfileUpdateContact{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateContactsBind.With(err))
	}

	// валидация модели
	if err := controller.validateUpdateContacts(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateContactsValidate.With(err))
	}

	if err := controller._updateContacts(&model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	UploadAvatar REST обертка для _uploadAvatar
*/
func (controller *Controller) UploadAvatarREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	// привязка модели
	model := models.ImageAvatar{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileAvatarBind.With(err))
	}

	// валидация модели
	if err := controller.validateUploadAvatar(&model); err != nil {
		return controller.Error(ctx, errors.ProfileAvatarValidate.With(err))
	}

	// загрузка аватара
	filePath, err := controller._uploadAvatar(session, token, &model)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, filePath)
}

/*
	UploadWallpaper Обертка для _uploadWallpaper
*/
func (controller *Controller) UploadWallpaperREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	// привязка модели
	model := models.ImageWallpaper{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileWallpaperBind.With(err))
	}

	// валидация модели
	if err := controller.validateUploadWallpaper(&model); err != nil {
		return controller.Error(ctx, errors.ProfileWallpaperValidate.With(err))
	}

	// загрузка фона
	filePath, err := controller._uploadWallpaper(session, token, &model)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, filePath)
}
