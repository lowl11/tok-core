package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение токена и логина
		token := ctx.Get("token").(string)
		username := ctx.Get("username").(string)

		// получить сессию с redis
		session, err := clientSession.Get(token, username)
		if err != nil {
			errorObject := &models.Response{
				Status:       "ERROR",
				Message:      "Произошла ошибка",
				InnerMessage: "Get session error",
			}
			return ctx.JSON(http.StatusUnauthorized, errorObject)
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
				errorObject := &models.Response{
					Status:       "ERROR",
					Message:      "Произошла ошибка",
					InnerMessage: "Get session error",
				}
				return ctx.JSON(http.StatusUnauthorized, errorObject)
			}
		} else if token != "" && username == "" {
			session, err = clientSession.GetByToken(token)
			if err != nil {
				errorObject := &models.Response{
					Status:       "ERROR",
					Message:      "Произошла ошибка",
					InnerMessage: "Get session error",
				}
				return ctx.JSON(http.StatusUnauthorized, errorObject)
			}
		} else if username != "" && token == "" {
			session, _, err = clientSession.GetByUsername(username)
			if err != nil {
				errorObject := &models.Response{
					Status:       "ERROR",
					Message:      "Произошла ошибка",
					InnerMessage: "Get session error",
				}
				return ctx.JSON(http.StatusUnauthorized, errorObject)
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
