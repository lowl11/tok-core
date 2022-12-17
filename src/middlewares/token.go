package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
)

func Token(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение и проверка токена
		token := ctx.Request().Header.Get("token")
		if token == "" {
			ctx.Error(errors.New("token is required"))
			return nil
		}

		// чтение и проверка логина
		username := ctx.Request().Header.Get("username")
		if token == "" {
			ctx.Error(errors.New("username is required"))
			return nil
		}

		// записать их в контекст
		ctx.Set("token", token)
		ctx.Set("username", username)

		if err := next(ctx); err != nil {
			ctx.Error(err)
			return nil
		}

		return nil
	}
}
