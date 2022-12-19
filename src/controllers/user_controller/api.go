package user_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

/*
	Info данные по пользователю
	Основные данные и подписки
*/
func (controller *Controller) Info(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// прочитать параметр
	username := ctx.Param("username")

	// получить сессию заданного пользователя
	userSession, _, err := controller.clientSession.GetByUsername(username)
	if err != nil && err.Error() != "session not found" {
		logger.Error(err, "Get user session error")
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	// у другого пользователя нет сессии, значит достаем из БД
	if userSession == nil {
		// получаем пользователя по юзернейму
		user, err := controller.userRepo.GetByUsername(username)
		if err != nil {
			logger.Error(err, "Get user by username error")
			return controller.Error(ctx, errors.UserGet.With(err))
		}

		// ошибка если пользователь не найден
		if user == nil {
			return controller.Error(ctx, errors.UserNotFound)
		}

		// получаем подписки пользователя
		subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(username)
		if err != nil {
			logger.Error(err, "Get user subscriptions error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		// получаем подписанных на пользователя
		subscribers, err := controller.subscriptRepo.ProfileSubscribers(username)
		if err != nil {
			logger.Error(err, "Get subscribers error")
			return controller.Error(ctx, errors.SubscribersGet.With(err))
		}

		// подписан ли тот чья сессия запросила
		var mySubscription bool
		if session != nil {
			mySubscription = subscribers.Single(func(item entities.ProfileSubscriber) bool {
				return item.Username == session.Username
			}) != nil
		}

		// обработка данных для клиента
		return controller.Ok(ctx, &models.UserInfoGet{
			MySubscription: mySubscription,

			Name: user.Name,
			BIO:  user.BIO,

			Avatar:    user.Avatar,
			Wallpaper: user.Wallpaper,

			Subscriptions: entities.ClientSessionSubscribes{
				SubscriberCount:   subscribers.Size(),
				SubscriptionCount: subscriptions.Size(),

				Subscriptions: subscriptions.Select(func(item entities.ProfileSubscription) string {
					return item.Username
				}).Slice(),
				Subscribers: subscribers.Select(func(item entities.ProfileSubscriber) string {
					return item.Username
				}).Slice(),
			},
		})
	}

	// получаем подписки пользователя
	subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(username)
	if err != nil {
		logger.Error(err, "Get user subscriptions error")
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	// получаем подписанных на пользователя
	subscribers, err := controller.subscriptRepo.ProfileSubscribers(username)
	if err != nil {
		logger.Error(err, "Get subscribers error")
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	// подписан ли тот чья сессия запросила
	var mySubscription bool
	if session != nil {
		mySubscription = subscribers.Single(func(item entities.ProfileSubscriber) bool {
			return item.Username == session.Username
		}) != nil
	}

	// обработка данных для клиента
	return controller.Ok(ctx, models.UserInfoGet{
		MySubscription: mySubscription,

		Name: userSession.Name,
		BIO:  userSession.BIO,

		Avatar:    userSession.Avatar,
		Wallpaper: userSession.Wallpaper,

		Subscriptions: entities.ClientSessionSubscribes{
			SubscriberCount:   subscribers.Size(),
			SubscriptionCount: subscriptions.Size(),

			Subscriptions: subscriptions.Select(func(item entities.ProfileSubscription) string {
				return item.Username
			}).Slice(),
			Subscribers: subscribers.Select(func(item entities.ProfileSubscriber) string {
				return item.Username
			}).Slice(),
		},
	})
}

/*
	Subscriptions возвращает список подписок (на кого)
*/
func (controller *Controller) Subscriptions(ctx echo.Context) error {
	logger := definition.Logger

	session := ctx.Get("client_session").(*entities.ClientSession)
	username := ctx.Param("username")

	// получаем список подписок (на кого)
	subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(username)
	if err != nil {
		logger.Error(err, "Get subscriptions error")
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	// для проверки подписана ли сессия
	mySubscriptions := array.NewWithList[string](session.Subscriptions.Subscriptions...)

	// обработка данных для клиента
	list := make([]models.UserSubscriptions, 0, subscriptions.Size())
	for subscriptions.Next() {
		item := subscriptions.Value()
		list = append(list, models.UserSubscriptions{
			Username:   item.Username,
			Name:       item.Name,
			Avatar:     item.Avatar,
			Subscribed: mySubscriptions.Contains(item.Username),
		})
	}

	return controller.Ok(ctx, list)
}

/*
	Subscribers возвращает список подписок (подписчики)
*/
func (controller *Controller) Subscribers(ctx echo.Context) error {
	logger := definition.Logger

	session := ctx.Get("client_session").(*entities.ClientSession)
	username := ctx.Param("username")

	// получаем подписчиков
	subscribers, err := controller.subscriptRepo.ProfileSubscribers(username)
	if err != nil {
		logger.Error(err, "Get subscribers error")
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	// проверяем подписана ли сессия
	mySubscriptions := array.NewWithList[string](session.Subscriptions.Subscriptions...)

	// обработка данных для клиента
	list := make([]models.UserSubscriber, 0, subscribers.Size())
	for subscribers.Next() {
		item := subscribers.Value()
		list = append(list, models.UserSubscriber{
			Username:   item.Username,
			Name:       item.Name,
			Avatar:     item.Avatar,
			Subscribed: mySubscriptions.Contains(item.Username),
		})
	}

	return controller.Ok(ctx, list)
}
