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

		// записать ее в контекст
		ctx.Set("token", token)

		if err := next(ctx); err != nil {
			ctx.Error(err)
			return nil
		}

		return nil
	}
}
