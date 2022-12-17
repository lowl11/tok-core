package profile_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) Subscribe(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	model := models.ProfileSubscribe{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileBind.With(err))
	}

	if err := controller.validateSubscribeProfile(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileValidate.With(err))
	}

	// сохранить подписку в БД
	if err := controller.subscriptRepo.ProfileSubscribe(session.Username, model.Username); err != nil {
		logger.Error(err, "Profile subscribe error")
		return controller.Error(ctx, errors.SubscribeOfProfile.With(err))
	}

	// обновить сессию в профиле
	session.Subscriptions.Subscriptions = append(session.Subscriptions.Subscriptions, model.Username)

	// обновить сессию другому пользователю на которого профиль подписался
	anotherSession, err := controller.clientSession.GetByUsername(model.Username)
	if err != nil && err.Error() != "session not found" {
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		anotherSession.Subscriptions.Subscribers = append(anotherSession.Subscriptions.Subscribers, session.Username)
	}

	// обновляем сессию профиля
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update session error")
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		// обновляем сессию другого на кого подписались
		if err = controller.clientSession.Update(anotherSession, token); err != nil {
			logger.Error(err, "Update session error")
			return controller.Error(ctx, errors.SessionUpdate.With(err))
		}
	}

	return controller.Ok(ctx, "OK")
}

func (controller *Controller) Unsubscribe(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	model := models.ProfileSubscribe{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileBind.With(err))
	}

	if err := controller.validateSubscribeProfile(&model); err != nil {
		return controller.Error(ctx, errors.SubscribeOfProfileValidate.With(err))
	}

	// сохранить подписку в БД
	if err := controller.subscriptRepo.ProfileUnsubscribe(session.Username, model.Username); err != nil {
		logger.Error(err, "Profile unsubscribe error")
		return controller.Error(ctx, errors.UnsubscribeOfProfile.With(err))
	}

	// обновить сессию в профиле
	session.Subscriptions.Subscriptions = append(session.Subscriptions.Subscriptions, model.Username)

	// обновить сессию другому пользователю на которого профиль подписался
	anotherSession, err := controller.clientSession.GetByUsername(model.Username)
	if err != nil && err.Error() != "session not found" {
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		anotherSession.Subscriptions.Subscribers = append(anotherSession.Subscriptions.Subscribers, session.Username)
	}

	// обновляем сессию профиля
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update session error")
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		// обновляем сессию другого на кого подписались
		if err = controller.clientSession.Update(anotherSession, token); err != nil {
			logger.Error(err, "Update session error")
			return controller.Error(ctx, errors.SessionUpdate.With(err))
		}
	}

	return controller.Ok(ctx, "OK")
}

func (controller *Controller) Update(ctx echo.Context) error {
	logger := definition.Logger
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

	// изменение данных пользователя
	if err := controller.userRepo.UpdateProfile(session.Username, &model); err != nil {
		logger.Error(err, "Update profile info error")
		return controller.Error(ctx, errors.ProfileUpdate.With(err))
	}

	// изменить данные из сессии
	session.BIO = &model.BIO
	session.Name = &model.Name

	if err := controller.clientSession.Update(session, token); err != nil {
		return controller.Error(ctx, errors.SessionCreate.With(err))
	}

	return controller.Ok(ctx, "OK")
}

func (controller *Controller) UploadAvatar(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	model := models.ImageAvatar{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileAvatarBind.With(err))
	}

	if err := controller.validateUploadAvatar(&model); err != nil {
		return controller.Error(ctx, errors.ProfileAvatarValidate.With(err))
	}

	fileName, err := controller.image.UploadAvatar(&model, session.Username)
	if err != nil {
		logger.Error(err, "Upload profile avatar error")
		return controller.Error(ctx, errors.ProfileAvatar.With(err))
	}

	if err = controller.userRepo.UpdateAvatar(session.Username, fileName); err != nil {
		return controller.Error(ctx, errors.ProfileUpdate.With(err))
	}

	filePath := "/images/profile/" + session.Username + "/" + fileName
	session.Avatar = &filePath
	if err = controller.clientSession.Update(session, token); err != nil {
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	return controller.Ok(ctx, filePath)
}

func (controller *Controller) UploadWallpaper(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)
	token := ctx.Get("token").(string)

	model := models.ImageWallpaper{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileWallpaperBind.With(err))
	}

	if err := controller.validateUploadWallpaper(&model); err != nil {
		return controller.Error(ctx, errors.ProfileWallpaperValidate.With(err))
	}

	fileName, err := controller.image.UploadWallpaper(&model, session.Username)
	if err != nil {
		logger.Error(err, "Upload profile wallpaper error")
		return controller.Error(ctx, errors.ProfileWallpaper.With(err))
	}

	if err = controller.userRepo.UpdateWallpaper(session.Username, fileName); err != nil {
		return controller.Error(ctx, errors.ProfileUpdate.With(err))
	}

	filePath := "/images/profile/" + session.Username + "/" + fileName
	session.Wallpaper = &filePath
	if err = controller.clientSession.Update(session, token); err != nil {
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	return controller.Ok(ctx, filePath)
}
