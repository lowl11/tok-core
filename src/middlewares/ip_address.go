package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
)

func IpAddress(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// чтение и проверка ip адреса
		ipAddress := ctx.RealIP()

		fmt.Println("remote:", ctx.Request().RemoteAddr)
		fmt.Println("x forwarded for:", ctx.Request().Header.Get("X-Forwarded-For"))
		fmt.Println("real ip:", ctx.Request().Header.Get("X-Real-Ip"))
		fmt.Println("")

		requestInJson, _ := json.Marshal(ctx.Request())
		fmt.Printf("%s\n", requestInJson)

		//if strings.Contains(ipAddress, "127.0.0.1") ||
		//	strings.Contains(ipAddress, "localhost") ||
		//	strings.Contains(ipAddress, "::1") {
		//	ipAddress = ctx.Request().RemoteAddr
		//}

		ctx.Set("ip_address", ipAddress)

		if err := next(ctx); err != nil {
			ctx.Error(err)
			return nil
		}

		return nil
	}
}
