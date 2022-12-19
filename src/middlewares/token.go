package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tok-core/src/data/models"
)

func Token(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение и проверка токена
		token := ctx.Request().Header.Get("token")
		if token == "" {
			errorObject := &models.Response{
				Status:       "ERROR",
				Message:      "Произошла ошибка",
				InnerMessage: "Token is required",
			}
			return ctx.JSON(http.StatusUnauthorized, errorObject)
		}

		// чтение и проверка логина
		username := ctx.Request().Header.Get("username")
		if token == "" {
			errorObject := &models.Response{
				Status:       "ERROR",
				Message:      "Произошла ошибка",
				InnerMessage: "Token is required",
			}
			return ctx.JSON(http.StatusUnauthorized, errorObject)
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

func TokenPart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение и проверка токена
		token := ctx.Request().Header.Get("token")

		// чтение и проверка логина
		username := ctx.Request().Header.Get("username")

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
