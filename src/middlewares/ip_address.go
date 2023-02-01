package middlewares

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/definition"
)

func IpAddress(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение и проверка ip адреса
		ipAddress := ctx.RealIP()
		definition.Logger.Warn("IP address from middleware: " + ipAddress)

		ctx.Set("ip_address", ipAddress)

		if err := next(ctx); err != nil {
			ctx.Error(err)
			return nil
		}

		return nil
	}
}
