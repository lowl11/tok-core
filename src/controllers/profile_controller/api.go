package profile_controller

import (
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
	"tok-core/src/services/action_helper"
)

/*
	_subscribe подписка профиля на "кого-то"
	Проверка на существование подписки
	Создается, запись в БД "кто на кого"
	Обновление сессии у того кто подписался и тому на кого подписались
*/
func (controller *Controller) _subscribe(session *entities.ClientSession, token string, model *models.ProfileSubscribe) *models.Error {
	logger := definition.Logger

	// проверить существует ли подписка
	exist, err := controller.subscriptRepo.Exist(session.Username, model.Username)
	if err != nil {
		return errors.SubscribersExist.With(err)
	}

	// если подписка уже существует
	if exist {
		return errors.SubscriptionExist
	}

	// сохранить подписку в БД
	if err = controller.subscriptRepo.ProfileSubscribe(session.Username, model.Username); err != nil {
		logger.Error(err, "Profile subscribe error")
		return errors.SubscribeOfProfile.With(err)
	}

	// обновить сессию в профиле
	session.Subscriptions.Subscriptions = append(session.Subscriptions.Subscriptions, model.Username)

	// обновить сессию другому пользователю на которого профиль подписался
	anotherSession, anotherToken, err := controller.clientSession.GetByUsername(model.Username)
	if err != nil && err.Error() != "session not found" {
		return errors.SessionGet.With(err)
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		anotherSession.Subscriptions.Subscribers = append(anotherSession.Subscriptions.Subscribers, session.Username)
	}

	// обновляем сессию профиля
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update session error")
		return errors.SessionUpdate.With(err)
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		// обновляем сессию другого на кого подписались
		if err = controller.clientSession.Update(anotherSession, anotherToken); err != nil {
			logger.Error(err, "Update session error")
			return errors.SessionUpdate.With(err)
		}
	}

	go func() {
		if err = controller.notification.Push(action_helper.Subscribe, anotherSession.Username, session.Username, nil); err != nil {
			logger.Error(err, "Push notification error", "Notifications")
		}
	}()

	return nil
}

/*
	_unsubscribe отписка от пользователя
	Проверка на существование подписки
	Удаление записи в БД "кто на кого"
	Обновление сессии у того кто отписался и тому от кого отписались
*/
func (controller *Controller) _unsubscribe(session *entities.ClientSession, token string, model *models.ProfileUnsubscribe) *models.Error {
	logger := definition.Logger

	// проверка на существование подписки
	exist, err := controller.subscriptRepo.Exist(session.Username, model.Username)
	if err != nil {
		return errors.SubscribersExist.With(err)
	}

	// ошибка если не существует
	if !exist {
		return errors.SubscriptionNotExist
	}

	// удалить подписку из БД
	if err = controller.subscriptRepo.ProfileUnsubscribe(session.Username, model.Username); err != nil {
		logger.Error(err, "Profile unsubscribe error")
		return errors.UnsubscribeOfProfile.With(err)
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
		return errors.SessionGet.With(err)
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
		return errors.SessionUpdate.With(err)
	}

	// вдруг у другого пользователя нет сессии
	if anotherSession != nil {
		// обновляем сессию другого на кого подписались
		if err = controller.clientSession.Update(anotherSession, anotherToken); err != nil {
			logger.Error(err, "Update session error")
			return errors.SessionUpdate.With(err)
		}
	}

	return nil
}

/*
	_update изменение "Имени" и "Био"
	Изменяется запись в БД
	Обновляется, сессия
*/
func (controller *Controller) _update(session *entities.ClientSession, token string, model *models.ProfileUpdate) *models.Error {
	logger := definition.Logger

	// изменение данных пользователя
	if err := controller.userRepo.UpdateProfile(session.Username, model); err != nil {
		logger.Error(err, "Update profile info error")
		return errors.ProfileUpdate.With(err)
	}

	// изменить данные из сессии
	session.BIO = &model.BIO
	session.Name = &model.Name

	// обновление сессии
	if err := controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update name & bio in session error")
		return errors.SessionCreate.With(err)
	}

	return nil
}

/*
	_updateContacts Обновление контактов профиля
	Обновляет, Телефон и Почту у пользователя
*/
func (controller *Controller) _updateContacts(model *models.ProfileUpdateContact) *models.Error {
	logger := definition.Logger

	// обновление контактов юзера в бд
	if err := controller.userRepo.UpdateContact(model); err != nil {
		logger.Error(err, "Update profile contacts error")
		return errors.ProfileUpdateContacts.With(err)
	}

	return nil
}

/*
	_uploadAvatar Загрузка аватара профиля
	Есть, валидация на расширения
	Удаляются, соседние файлы если они есть
	Обновляется, сессия
*/
func (controller *Controller) _uploadAvatar(session *entities.ClientSession, token string, model *models.ImageAvatar) (string, *models.Error) {
	logger := definition.Logger

	// загрузка изображения
	fileName, err := controller.image.UploadAvatar(model, session.Username)
	if err != nil {
		logger.Error(err, "Upload profile avatar error")
		return "", errors.ProfileAvatar.With(err)
	}

	// обновление пути к аватару в БД
	if err = controller.userRepo.UpdateAvatar(session.Username, fileName); err != nil {
		logger.Error(err, "Update avatar is DB error")
		return "", errors.ProfileUpdate.With(err)
	}

	// обновление пути к аватару в сессии пользователя
	filePath := "/images/profile/" + session.Username + "/" + fileName
	session.Avatar = &filePath
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update avatar in session")
		return "", errors.SessionUpdate.With(err)
	}

	return filePath, nil
}

/*
	_uploadWallpaper Загрузка фона профиля
	Есть, валидация на расширения
	Удаляются, соседние файлы если они есть
	Обновляется, сессия
*/
func (controller *Controller) _uploadWallpaper(session *entities.ClientSession, token string, model *models.ImageWallpaper) (string, *models.Error) {
	logger := definition.Logger

	// загрузка фона
	fileName, err := controller.image.UploadWallpaper(model, session.Username)
	if err != nil {
		logger.Error(err, "Upload profile wallpaper error")
		return "", errors.ProfileWallpaper.With(err)
	}

	// обновление пути к фону в БД
	if err = controller.userRepo.UpdateWallpaper(session.Username, fileName); err != nil {
		logger.Error(err, "Update wallpaper in DB error")
		return "", errors.ProfileUpdate.With(err)
	}

	// обновление пути к фону в сессии пользователя
	filePath := "/images/profile/" + session.Username + "/" + fileName
	session.Wallpaper = &filePath
	if err = controller.clientSession.Update(session, token); err != nil {
		logger.Error(err, "Update wallpaper in session error")
		return "", errors.SessionUpdate.With(err)
	}

	return filePath, nil
}
