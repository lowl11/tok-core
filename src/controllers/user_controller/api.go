package user_controller

import (
	"github.com/labstack/echo/v4"
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

	return controller.Ok(ctx, &models.ClientSessionGet{
		Username: user.Username,
		Name:     user.Name,
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

	session, err := controller.clientSession.Get(model.Token)
	if err != nil {
		logger.Error(err, "Get client session error")
		return controller.Error(ctx, errors.SessionGet.With(err))
	}

	if session == nil {
		return controller.Error(ctx, errors.SessionNotFound)
	}

	return controller.Ok(ctx, &models.ClientSessionGet{
		Username: session.Username,
		Name:     session.Name,
	})
}
