package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
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

func SessionPart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение токена и логина
		token := ctx.Get("token").(string)
		username := ctx.Get("username").(string)

		var session *entities.ClientSession
		var err error

		// получить сессию с redis
		if token != "" && username != "" {
			session, err = clientSession.Get(token, username)
			if err != nil {
				ctx.Error(errors.New("get session error"))
				return nil
			}
		} else if token != "" && username == "" {
			session, err = clientSession.GetByToken(token)
			if err != nil {
				ctx.Error(errors.New("get session error"))
				return nil
			}
		} else if username != "" && token == "" {
			session, err = clientSession.GetByUsername(username)
			if err != nil {
				ctx.Error(errors.New("get session error"))
				return nil
			}
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
