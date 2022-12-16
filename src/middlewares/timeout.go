package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

var (
	timeout = time.Second * 5
)

func Timeout() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Запрос достигнул таймаута",
		Timeout:      timeout,

		OnTimeoutRouteErrorHandler: func(err error, ctx echo.Context) {
			timeoutError := errors.New("request reached timeout | " + err.Error())
			ctx.Error(timeoutError)
		},
	})
}