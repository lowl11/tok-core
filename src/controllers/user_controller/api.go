package user_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) Info(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// прочитать параметр
	username := ctx.Param("username")

	// получить сессию заданного пользователя
	userSession, err := controller.clientSession.GetByUsername(username)
	if err != nil && err.Error() != "session not found" {
		logger.Error(err, "Get user session error")
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	// у другого пользователя нет сессии, значит достаем из БД
	if userSession == nil {
		user, err := controller.userRepo.GetByUsername(username)
		if err != nil {
			logger.Error(err, "Get user by username error")
			return controller.Error(ctx, errors.UserGet.With(err))
		}

		if user == nil {
			return controller.Error(ctx, errors.UserNotFound)
		}

		subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(username)
		if err != nil {
			logger.Error(err, "Get user subscriptions error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		subscribers, err := controller.subscriptRepo.ProfileSubscribers(username)
		if err != nil {
			logger.Error(err, "Get subscribers error")
			return controller.Error(ctx, errors.SubscribersGet.With(err))
		}

		return controller.Ok(ctx, &models.UserInfoGet{
			MySubscription: array.NewWithList[string](session.Subscriptions.Subscriptions...).ContainsFunc(func(item string) bool {
				return item == username
			}),

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

	subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(username)
	if err != nil {
		logger.Error(err, "Get user subscriptions error")
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	subscribers, err := controller.subscriptRepo.ProfileSubscribers(username)
	if err != nil {
		logger.Error(err, "Get subscribers error")
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	return controller.Ok(ctx, models.UserInfoGet{
		MySubscription: array.NewWithList[string](session.Subscriptions.Subscriptions...).ContainsFunc(func(item string) bool {
			return item == username
		}),

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

func (controller *Controller) Subscriptions(ctx echo.Context) error {
	logger := definition.Logger

	username := ctx.Param("username")

	subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(username)
	if err != nil {
		logger.Error(err, "Get subscriptions error")
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	list := make([]models.UserSubscriptions, 0, subscriptions.Size())
	for subscriptions.Next() {
		item := subscriptions.Value()
		list = append(list, models.UserSubscriptions{
			Username:   item.Username,
			Name:       item.Name,
			Avatar:     item.Avatar,
			Subscribed: false, // TODO
		})
	}

	return controller.Ok(ctx, list)
}

func (controller *Controller) Subscribers(ctx echo.Context) error {
	logger := definition.Logger

	username := ctx.Param("username")

	subscribers, err := controller.subscriptRepo.ProfileSubscribers(username)
	if err != nil {
		logger.Error(err, "Get subscribers error")
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	list := make([]models.UserSubscriber, 0, subscribers.Size())
	for subscribers.Next() {
		item := subscribers.Value()
		list = append(list, models.UserSubscriber{
			Username:   item.Username,
			Name:       item.Name,
			Avatar:     item.Avatar,
			Subscribed: false, // TODO
		})
	}

	return controller.Ok(ctx, list)
}

func (controller *Controller) Search(ctx echo.Context) error {
	// query
	return controller.Ok(ctx, "OK")
}
