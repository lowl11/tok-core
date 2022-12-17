package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
)

func Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение токена и логина
		token := ctx.Get("token").(string)
		username := ctx.Get("username").(string)

		// получить сессию с redis
		session, err := clientSession.Get(token, username)
		if err != nil {
			ctx.Error(errors.New("get session error"))
			return nil
		}

		// записать ее в контекст
		ctx.Set("client_session", session)

		if err = next(ctx); err != nil {
			ctx.Error(err)
			return nil
		}

		return nil
	}
}
