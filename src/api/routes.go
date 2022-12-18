package api

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/controllers"
	"tok-core/src/controllers/auth_controller"
	"tok-core/src/controllers/feed_controller"
	"tok-core/src/controllers/post_controller"
	"tok-core/src/controllers/profile_controller"
	"tok-core/src/controllers/user_controller"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/middlewares"
	"tok-core/src/repositories"
)

func setRoutes(server *echo.Echo) {
	logger := definition.Logger

	// ивенты
	apiEvents, err := events.Get()
	if err != nil {
		logger.Fatal(err, "Initializing events error")
	}

	// проставить клиентские сессии глобально
	middlewares.SetClientSession(apiEvents.ClientSession)

	// репозитории
	apiRepositories, err := repositories.Get(apiEvents)
	if err != nil {
		logger.Fatal(err, "Connecting to database error")
	}

	// контроллеры
	apiControllers := controllers.Get(apiRepositories, apiEvents)

	// статичные методы
	server.GET("/health", apiControllers.Static.Health)
	server.RouteNotFound("*", apiControllers.Static.RouteNotFound)

	// эндпоинты
	setAuth(server, apiControllers.Auth)
	setProfile(server, apiControllers.Profile)
	setUser(server, apiControllers.User)

	setFeed(server, apiControllers.Feed)
	setPost(server, apiControllers.Post)
}

func setAuth(server *echo.Echo, controller *auth_controller.Controller) {
	group := server.Group("/api/v1/auth")

	group.POST("/signup", controller.Signup)

	group.POST("/login/credentials", controller.LoginByCredentials)
	group.POST("/login/token", controller.LoginByToken)
}

func setProfile(server *echo.Echo, controller *profile_controller.Controller) {
	group := server.Group("/api/v1/profile")

	setMiddlewaresPublic(group)
	group.POST("/update", controller.Update)
	group.POST("/avatar", controller.UploadAvatar)
	group.POST("/wallpaper", controller.UploadWallpaper)

	group.POST("/subscribe", controller.Subscribe)
	group.POST("/unsubscribe", controller.Unsubscribe)
}

func setUser(server *echo.Echo, controller *user_controller.Controller) {
	// part access group
	partGroup := server.Group("/api/v1/user")
	setMiddlewaresPublicPart(partGroup)

	partGroup.GET("/info/:username", controller.Info)

	// full access group

	group := server.Group("/api/v1/user")

	setMiddlewaresPublic(group)
	group.GET("/subscribers/:username", controller.Subscribers)
	group.GET("/subscriptions/:username", controller.Subscriptions)
	group.GET("/search", controller.Search)
}

func setFeed(server *echo.Echo, controller *feed_controller.Controller) {
	group := server.Group("/api/v1/feed")

	setMiddlewaresPublic(group)
}

func setPost(server *echo.Echo, controller *post_controller.Controller) {
	group := server.Group("/api/v1/post")

	setMiddlewaresPublic(group)
}
