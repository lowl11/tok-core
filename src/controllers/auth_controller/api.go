package auth_controller

import (
	"github.com/labstack/echo/v4"
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
	"tok-core/src/services/string_encryptor"
)

/*
	Signup регистрация новой учетной записи
	После заведения пользователя в БД сразу создает сессию и позволяет авторизоваться по выданному токену
*/
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

	// создать модель сессии
	sessionModel := &models.ClientSessionCreate{
		Username: model.Username,
	}

	// создать сессию в Redis
	var userToken string
	if userToken, err = controller.clientSession.Create(sessionModel); err != nil {
		logger.Error(err, "Create user session error")
		return controller.Error(ctx, errors.SessionCreate.With(err))
	}

	return controller.Ok(ctx, userToken)
}

/*
	LoginByCredentials авторизация с юзернеймом и паролем
	Возвращает: сгенерированный токен сессии по которому можно будет в дальнейшем авторизоваться
*/
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

	model.Username = controller.preprocessUsername(model.Username)

	// получить пользователя
	user, err := controller.userRepo.GetByUsername(model.Username)
	if err != nil {
		logger.Error(err, "Get user error", layers.Auth)
		return controller.Unauthorized(ctx, errors.UserGet.With(err))
	}

	// если пользователь не найден
	if user == nil {
		return controller.Unauthorized(ctx, errors.UserNotFound)
	}

	// расшифровать пароль
	decryptedPassword, err := string_encryptor.Decrypt(config.User.CryptKey, user.Password)
	if err != nil {
		logger.Error(err, "Decrypting password error")
		return controller.Error(ctx, errors.UserDecryptPassword)
	}

	// сверка паролей
	if model.Password != decryptedPassword {
		return controller.Error(ctx, errors.LoginPassword)
	}

	// проверка существует ли уже сессия у пользователя
	session, token, err := controller.clientSession.GetByUsername(model.Username)
	if err != nil && err.Error() != "session not found" {
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	if session != nil {
		return controller.Ok(ctx, &models.ClientSessionGet{
			Token:    token,
			Username: user.Username,

			Name: user.Name,
			BIO:  user.BIO,

			Avatar:    user.Avatar,
			Wallpaper: user.Wallpaper,

			Subscriptions: entities.ClientSessionSubscribes{
				SubscriberCount:   session.Subscriptions.SubscriberCount,
				SubscriptionCount: session.Subscriptions.SubscriptionCount,

				Subscribers:   session.Subscriptions.Subscribers,
				Subscriptions: session.Subscriptions.Subscriptions,
			},
		})
	}

	// удаление всех сессий с таким логином
	if err = controller.clientSession.DeleteByUsername(model.Username); err != nil {
		return controller.Error(ctx, errors.SessionDelete.With(err))
	}

	// получить список подписчиков сессии
	subscribers, err := controller.subscriptRepo.ProfileSubscribers(model.Username)
	if err != nil {
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	// получить список подписок сессии
	subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(model.Username)
	if err != nil {
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	// создание сессии
	var sessionToken string
	if sessionToken, err = controller.clientSession.Create(&models.ClientSessionCreate{
		Username:  model.Username,
		Name:      user.Name,
		BIO:       user.BIO,
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
	}); err != nil {
		return controller.Error(ctx, errors.SessionCreate.With(err))
	}

	// проверить нужно ли запомнить ip адрес
	ipAddress := ctx.Get("ip_address").(string)
	if model.Remember && ipAddress != "" {
		logger.Info("Login by credentials. IP address: " + ipAddress)
		go func() {
			logger.Info("Going to save user IP address", layers.Auth)

			// проверка существует ли уже ip адрес
			userIp, err := controller.userIpRepo.GetByIpAddress(ipAddress)
			if err != nil {
				logger.Error(err, "Get user by ip address error", layers.Database)
				return
			}

			// если запись найдена и ip принадлежит другому пользователю
			if userIp != nil {
				if userIp.Username != model.Username {
					if err = controller.userIpRepo.DeleteByIP(userIp.IpAddress); err != nil {
						logger.Error(err, "", layers.Database)
						return
					}
				} else {
					// если запись уже найдена, просто идем дальше
					return
				}
			}

			if err = controller.userIpRepo.New(model.Username, ipAddress); err != nil {
				logger.Error(err, "Creating new bind username + ip address error", layers.Database)
				return
			}
		}()
	}

	// обработка данных для клиента
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

/*
	LoginByToken авторизация по токену сессии
*/
func (controller *Controller) LoginByToken(ctx echo.Context) error {
	logger := definition.Logger
	ipAddress := ctx.Get("ip_address").(string)

	// привязка модели
	model := models.LoginByToken{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.LoginBind.With(err))
	}

	// валидация модели
	if err := controller.validateLoginByToken(&model); err != nil {
		return controller.Error(ctx, errors.LoginValidate.With(err))
	}

	// если токен не найден, пытаемся авторизоваться через ip
	if model.Token == "" && ipAddress != "" {
		userIp, err := controller.userIpRepo.GetByIpAddress(ipAddress)
		if err != nil {
			logger.Error(err, "Get username by ip address error", layers.Database)
			return controller.Unauthorized(ctx, errors.UserIpGetByIp.With(err))
		}

		if userIp == nil {
			return controller.Unauthorized(ctx, errors.UserIpUserNotFound)
		}

		session, token, err := controller.clientSession.GetByUsername(userIp.Username)
		if err != nil {
			logger.Error(err, "Get user session error", layers.Redis)
			return controller.Unauthorized(ctx, errors.SessionGet.With(err))
		}

		// если сессию нашли, возвращаем ее
		if session != nil {
			return controller.Ok(ctx, &models.ClientSessionGet{
				Token:    token,
				Username: session.Username,

				Name: session.Name,
				BIO:  session.BIO,

				Avatar:    session.Avatar,
				Wallpaper: session.Wallpaper,

				Subscriptions: session.Subscriptions,
			})
		}

		userByIp, err := controller.userRepo.GetByUsername(userIp.Username)
		if err != nil {
			logger.Error(err, "Get user by username error", layers.Database)
			return controller.Error(ctx, errors.UserGet.With(err))
		}

		// если сессию не нашли, создадим ее
		newSessionToken, err := controller.clientSession.Create(&models.ClientSessionCreate{
			Username:  userByIp.Username,
			Name:      userByIp.Name,
			BIO:       userByIp.BIO,
			Avatar:    userByIp.Avatar,
			Wallpaper: userByIp.Wallpaper,
		})
		if err != nil {
			logger.Error(err, "Create session error", layers.Redis)
			return controller.Error(ctx, errors.SessionCreate.With(err))
		}

		// получение списка подписок сессии
		subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(session.Username)
		if err != nil {
			logger.Error(err, "Get profile subscriptions error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		// получение списка подписчиков сессии
		subscribers, err := controller.subscriptRepo.ProfileSubscribers(userByIp.Username)
		if err != nil {
			logger.Error(err, "Get profile subscribers error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		return controller.Ok(ctx, &models.ClientSessionGet{
			Token:    newSessionToken,
			Username: userByIp.Username,

			Name: userByIp.Name,
			BIO:  userByIp.BIO,

			Avatar:    userByIp.Avatar,
			Wallpaper: userByIp.Wallpaper,

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

	// получение сессии по токену
	session, err := controller.clientSession.GetByToken(model.Token)
	if err != nil {
		logger.Error(err, "Get client session error")
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	// если сессия не найдена
	if session == nil {
		return controller.Unauthorized(ctx, errors.SessionNotFound)
	}

	// кол-во подписок сессии
	subscriptionsCount, err := controller.subscriptRepo.ProfileSubscriptionsCount(session.Username)
	if err != nil {
		logger.Error(err, "Get profile subscriptions count error")
		return controller.Error(ctx, errors.SubscriptionsGet.With(err))
	}

	// кол-во подписчиков сессии
	subscribersCount, err := controller.subscriptRepo.ProfileSubscribersCount(session.Username)
	if err != nil {
		logger.Error(err, "Get profile subscribers count error")
		return controller.Error(ctx, errors.SubscribersGet.With(err))
	}

	// нужно ли обновлять сессию
	var sessionUpdate bool

	// если кол-во подписок (на кого) не совпадает
	if session.Subscriptions.SubscriptionCount != subscriptionsCount {
		// получаем список подписок сессии
		subscriptions, err := controller.subscriptRepo.ProfileSubscriptions(session.Username)
		if err != nil {
			logger.Error(err, "Get profile subscriptions error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		session.Subscriptions.Subscriptions = subscriptions.Select(func(item entities.ProfileSubscription) string {
			return item.Username
		}).Slice()
		session.Subscriptions.SubscriptionCount = len(session.Subscriptions.Subscriptions)
		sessionUpdate = true
	}

	// если кол-во подписчиков (кто на него) не совпадает
	if session.Subscriptions.SubscriberCount != subscribersCount {
		// получение списка подписчиков сессии
		subscribers, err := controller.subscriptRepo.ProfileSubscribers(session.Username)
		if err != nil {
			logger.Error(err, "Get profile subscribers error")
			return controller.Error(ctx, errors.SubscriptionsGet.With(err))
		}

		session.Subscriptions.Subscribers = subscribers.Select(func(item entities.ProfileSubscriber) string {
			return item.Username
		}).Slice()
		session.Subscriptions.SubscriberCount = len(session.Subscriptions.Subscribers)
		sessionUpdate = true
	}

	// проверяем нужно ли обновить сессию
	if sessionUpdate {
		// обновляем сессию
		if err = controller.clientSession.Update(session, model.Token); err != nil {
			logger.Error(err, "Session update error")
			return controller.Error(ctx, errors.SessionUpdate.With(err))
		}
	}

	// обработка данных для клиента
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
