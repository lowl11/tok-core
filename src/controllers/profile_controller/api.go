package profile_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

/*
	Subscribe подписка профиля на "кого-то"
	Проверка на существование подписки
	Создается запись в БД "кто на кого"
	Обнолвение сессии у того кто подписался и тому на кого подписались
*/
func (controller *Controller) Subscribe(ctx echo.Context) error {
	logger := definition.Logger
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

	// проверить существует ли подписка
	exist, err := controller.subscriptRepo.Exist(session.Username, model.Username)
	if err != nil {
		return controller.Error(ctx, errors.SubscribersExist.With(err))
	}

	// если подписка уже существует
	if exist {
		return controller.Error(ctx, errors.SubscriptionExist)
	}

	// сохранить подписку в БД
	if err = controller.subscriptRepo.ProfileSubscribe(session.Username, model.Username); err != nil {
		logger.Error(err, "Profile subscribe error")
		return controller.Error(ctx, errors.SubscribeOfProfile.With(err))
	}

	// обновить сессию в профиле
	session.Subscriptions.Subscriptions = append(session.Subscriptions.Subscriptions, model.Username)

	// обновить сессию другому пользователю на которого профиль подписался
	anotherSession, anotherToken, err := controller.clientSession.GetByUsername(model.Username)
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
		if err = controller.clientSession.Update(anotherSession, anotherToken); err != nil {
			logger.Error(err, "Update session error")
			return controller.Error(ctx, errors.SessionUpdate.With(err))
		}
	}

	return controller.Ok(ctx, "OK")
}

/*
	Unsubscribe отписка от пользователя
	Проверка на существование подписки
	Удаление записи в БД "кто на кого"
	Обнолвение сессии у того кто отписался и тому от кого отписались
*/
func (controller *Controller) Unsubscribe(ctx echo.Context) error {
	logger := definition.Logger
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

	// проверка на существование подписки
	exist, err := controller.subscriptRepo.Exist(session.Username, model.Username)
	if err != nil {
		return controller.Error(ctx, errors.SubscribersExist.With(err))
	}

	// ошибка если не существует
	if !exist {
		return controller.Error(ctx, errors.SubscriptionNotExist)
	}

	// удалить подписку из БД
	if err = controller.subscriptRepo.ProfileUnsubscribe(session.Username, model.Username); err != nil {
		logger.Error(err, "Profile unsubscribe error")
		return controller.Error(ctx, errors.UnsubscribeOfProfile.With(err))
	}

	// обновить сессию в профиле
	profileList := array.NewWithList[string](session.Subscriptions.Subscriptions...)
	removeIndex := profileList.IndexOf(model.Username)
	if removeIndex > -1 {
		profileList.Remove(removeIndex)
	}
	session.Subscriptions.Subscriptions = profileList.Slice()

	// обновить сессию другому пользователю на которого профиль подписался
	anotherSession, anotherToken, err := controller.clientSession.GetByUsername(model.Username)
	if err != nil && err.Error() != "session not found" {
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		userList := array.NewWithList[string](anotherSession.Subscriptions.Subscriptions...)
		userRemoveIndex := userList.IndexOf(anotherSession.Username)
		if userRemoveIndex > -1 {
			userList.Remove(userRemoveIndex)
		}
		anotherSession.Subscriptions.Subscribers = userList.Slice()
	}

	// обновляем сессию профиля
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update session error")
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		// обновляем сессию другого на кого подписались
		if err = controller.clientSession.Update(anotherSession, anotherToken); err != nil {
			logger.Error(err, "Update session error")
			return controller.Error(ctx, errors.SessionUpdate.With(err))
		}
	}

	return controller.Ok(ctx, "OK")
}

/*
	Update изменение "Имени" и "Био"
	Изменяется запись в БД
	Обновляется сессия
*/
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

	// обновление сессии
	if err := controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update name & bio in session error")
		return controller.Error(ctx, errors.SessionCreate.With(err))
	}

	return controller.Ok(ctx, "OK")
}

/*
	UploadAvatar Загрузка аватара профиля
	Есть валидация на расширения
	Удаляются соседние файлы если они есть
	Обновляется сессия
*/
func (controller *Controller) UploadAvatar(ctx echo.Context) error {
	logger := definition.Logger
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

	// загрузка изображения
	fileName, err := controller.image.UploadAvatar(&model, session.Username)
	if err != nil {
		logger.Error(err, "Upload profile avatar error")
		return controller.Error(ctx, errors.ProfileAvatar.With(err))
	}

	// обновление пути к аватару в БД
	if err = controller.userRepo.UpdateAvatar(session.Username, fileName); err != nil {
		logger.Error(err, "Update avatar is DB error")
		return controller.Error(ctx, errors.ProfileUpdate.With(err))
	}

	// обновление пути к аватару в сессии пользователя
	filePath := "/images/profile/" + session.Username + "/" + fileName
	session.Avatar = &filePath
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update avatar in session")
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	return controller.Ok(ctx, filePath)
}

/*
	UploadWallpaper Загрузка фона профиля
	Есть валидация на расширения
	Удаляются соседние файлы если они есть
	Обновляется сессия
*/
func (controller *Controller) UploadWallpaper(ctx echo.Context) error {
	logger := definition.Logger
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
	fileName, err := controller.image.UploadWallpaper(&model, session.Username)
	if err != nil {
		logger.Error(err, "Upload profile wallpaper error")
		return controller.Error(ctx, errors.ProfileWallpaper.With(err))
	}

	// обновление пути к фону в БД
	if err = controller.userRepo.UpdateWallpaper(session.Username, fileName); err != nil {
		logger.Error(err, "Update wallpaper in DB error")
		return controller.Error(ctx, errors.ProfileUpdate.With(err))
	}

	// обновление пути к фону в сессии пользователя
	filePath := "/images/profile/" + session.Username + "/" + fileName
	session.Wallpaper = &filePath
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update wallpaper in session error")
		return controller.Error(ctx, errors.SessionUpdate.With(err))
	}

	return controller.Ok(ctx, filePath)
}

func (controller *Controller) UpdateContacts(ctx echo.Context) error {
	logger := definition.Logger

	// привязка модели
	model := models.ProfileUpdateContact{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateContactsBind.With(err))
	}

	// валидация модели
	if err := controller.validateUpdateContacts(&model); err != nil {
		return controller.Error(ctx, errors.ProfileUpdateContactsValidate.With(err))
	}

	// обновление контактов юзера в бд
	if err := controller.userRepo.UpdateContact(&model); err != nil {
		logger.Error(err, "Update profile contacts error")
		return controller.Error(ctx, errors.ProfileUpdateContacts.With(err))
	}

	return controller.Ok(ctx, "OK")
}
