package auth_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
	"tok-core/src/services/string_encryptor"
)

func (controller *Controller) Signup(ctx echo.Context) error {
	logger := definition.Logger
	config := definition.Config

	// привязка модели
	model := models.Signup{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.SignupBind.With(err))
	}

	// обработка переменных
	model.Username = controller.preprocessUsername(model.Username)

	// валидация модели
	if err := controller.validateSignUp(&model); err != nil {
		return controller.Error(ctx, errors.SignupValidate.With(err))
	}

	// шифровка пароля
	encryptedPassword, err := string_encryptor.Encrypt(config.User.CryptKey, model.Password)
	if err != nil {
		logger.Error(err, "Encrypting password error")
		return controller.Error(ctx, errors.UserEncryptPassword.With(err))
	}

	// завести пользователя в БД
	if err = controller.authRepo.Signup(&model, encryptedPassword); err != nil {
		logger.Error(err, "Sign up user error")
		return controller.Error(ctx, errors.Signup.With(err))
	}

	// создать сессию в Redis
	sessionModel := &models.ClientSessionCreate{
		Username: model.Username,
	}

	var userToken string
	if userToken, err = controller.clientSession.Create(sessionModel); err != nil {
		logger.Error(err, "Create user session error")
		return controller.Error(ctx, errors.SessionCreate.With(err))
	}

	return controller.Ok(ctx, userToken)
}

func (controller *Controller) LoginByCredentials(ctx echo.Context) error {
	logger := definition.Logger
	config := definition.Config

	// привязка модели
	model := models.LoginByCredentials{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.LoginBind.With(err))
	}

	// валидация модели
	if err := controller.validateLoginByCredentials(&model); err != nil {
		return controller.Error(ctx, errors.LoginValidate.With(err))
	}

	// получить пользователя
	user, err := controller.userRepo.GetByUsername(model.Username)
	if err != nil {
		logger.Error(err, "Get user error")
		return controller.Error(ctx, errors.UserGet.With(err))
	}

	if user == nil {
		return controller.Error(ctx, errors.UserNotFound)
	}

	// расшифровать пароль
	decryptedPassword, err := string_encryptor.Decrypt(config.User.CryptKey, user.Password)
	if err != nil {
		logger.Error(err, "Decrypting password error")
		return controller.Error(ctx, errors.UserDecryptPassword)
	}

	// сверка паролей
	if model.Password != decryptedPassword {
		return controller.Error(ctx, errors.LoginPassword.With(err))
	}

	// удаление всех сессий с таким логином
	if err = controller.clientSession.DeleteByUsername(model.Username); err != nil {
		return controller.Error(ctx, errors.SessionDelete.With(err))
	}

	// создание сессии
	var sessionToken string
	if sessionToken, err = controller.clientSession.Create(&models.ClientSessionCreate{
		Username:  model.Username,
		Name:      user.Name,
		BIO:       user.BIO,
		Avatar:    user.Avatar,
		Wallpaper: user.Wallpaper,
	}); err != nil {
		return controller.Error(ctx, errors.SessionCreate.With(err))
	}

	// получить список подписок и подписчиков
	subscribers, err := controller.subscriptRepo.ProfileSubscribers(model.Username)
	if err != nil {
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(model.Username)
	if err != nil {
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	return controller.Ok(ctx, &models.ClientSessionGet{
		Token:    sessionToken,
		Username: user.Username,

		Name: user.Name,
		BIO:  user.BIO,

		Avatar:    user.Avatar,
		Wallpaper: user.Wallpaper,

		Subscriptions: entities.ClientSessionSubscribes{
			SubscriberCount:   subscribers.Size(),
			SubscriptionCount: subscriptions.Size(),

			Subscribers: subscribers.Select(func(item entities.ProfileSubscriber) string {
				return item.Username
			}).Slice(),
			Subscriptions: subscriptions.Select(func(item entities.ProfileSubscription) string {
				return item.Username
			}).Slice(),
		},
	})
}

func (controller *Controller) LoginByToken(ctx echo.Context) error {
	logger := definition.Logger

	// привязка модели
	model := models.LoginByToken{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.LoginBind.With(err))
	}

	// валидация модели
	if err := controller.validateLoginByToken(&model); err != nil {
		return controller.Error(ctx, errors.LoginValidate.With(err))
	}

	session, err := controller.clientSession.GetByToken(model.Token)
	if err != nil {
		logger.Error(err, "Get client session error")
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	if session == nil {
		return controller.Error(ctx, errors.SessionNotFound)
	}

	subscriptionsCount, err := controller.subscriptRepo.ProfileSubscriptionsCount(session.Username)
	if err != nil {
		logger.Error(err, "Get profile subscriptions count error")
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	subscribersCount, err := controller.subscriptRepo.ProfileSubscribersCount(session.Username)
	if err != nil {
		logger.Error(err, "Get profile subscribers count error")
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	// если кол-во подписок (на кого) не совпадает
	if len(session.Subscriptions.Subscriptions) != subscriptionsCount {
		subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(session.Username)
		if err != nil {
			logger.Error(err, "Get profile subscriptions error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		session.Subscriptions.Subscriptions = subscriptions.Select(func(item entities.ProfileSubscription) string {
			return item.Username
		}).Slice()
		session.Subscriptions.SubscriptionCount = len(session.Subscriptions.Subscriptions)
	}

	// если кол-во подписчиков (кто на него) не совпадает
	if len(session.Subscriptions.Subscribers) != subscribersCount {
		subscribers, err := controller.subscriptRepo.ProfileSubscribers(session.Username)
		if err != nil {
			logger.Error(err, "Get profile subscribers error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		session.Subscriptions.Subscribers = subscribers.Select(func(item entities.ProfileSubscriber) string {
			return item.Username
		}).Slice()
		session.Subscriptions.SubscriberCount = len(session.Subscriptions.Subscribers)
	}

	return controller.Ok(ctx, &models.ClientSessionGet{
		Token:    model.Token,
		Username: session.Username,

		Name: session.Name,
		BIO:  session.BIO,

		Avatar:    session.Avatar,
		Wallpaper: session.Wallpaper,

		Subscriptions: session.Subscriptions,
	})
}
