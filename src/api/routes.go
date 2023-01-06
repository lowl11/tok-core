package api

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/controllers"
	"tok-core/src/controllers/auth_controller"
	"tok-core/src/controllers/feed_controller"
	"tok-core/src/controllers/post_controller"
	"tok-core/src/controllers/profile_controller"
	"tok-core/src/controllers/search_controller"
	"tok-core/src/controllers/user_controller"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/middlewares"
	"tok-core/src/repositories"
	"tok-core/src/services/feed_helper"
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

	feed_helper.SetLikeRepository(apiRepositories.PostLike)
	feed_helper.SetCommentRepository(apiRepositories.PostComment)

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
	setSearch(server, apiControllers.Search)
}

func setAuth(server *echo.Echo, controller *auth_controller.Controller) {
	group := server.Group("/api/v1/auth")

	group.POST("/signup", controller.Signup)

	group.POST("/login/credentials", controller.LoginByCredentials)
	group.POST("/login/token", controller.LoginByToken)
	group.POST("/login/ip", controller.LoginByIP)
}

func setProfile(server *echo.Echo, controller *profile_controller.Controller) {
	group := server.Group("/api/v1/profile")

	setMiddlewaresPublic(group)
	group.POST("/update", controller.UpdateREST)
	group.POST("/avatar", controller.UploadAvatarREST)
	group.POST("/wallpaper", controller.UploadWallpaperREST)

	group.POST("/subscribe", controller.SubscribeREST)
	group.POST("/unsubscribe", controller.UnsubscribeREST)
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
}

func setFeed(server *echo.Echo, controller *feed_controller.Controller) {
	group := server.Group("/api/v1/feed")

	setMiddlewaresPublic(group)
	group.GET("/main", controller.MainREST)
	group.GET("/explore", controller.ExploreREST)
	group.GET("/user/:username", controller.UserREST)
	group.GET("/category/:category_code", controller.CategoryREST)
	group.GET("/single/:code", controller.SingleREST)
}

func setPost(server *echo.Echo, controller *post_controller.Controller) {
	group := server.Group("/api/v1/post")

	setMiddlewaresPublic(group)

	// main
	group.POST("/add", controller.AddREST)
	group.DELETE("/delete/:code", controller.DeleteREST)

	// likes
	group.POST("/like/do", controller.LikeREST)
	group.POST("/like/undo", controller.UnlikeREST)
	group.GET("/like/get/:code", controller.GetLikesREST)

	// comments
	group.POST("/comment/add", controller.AddCommentREST)
	group.GET("/comment/get", controller.GetCommentREST)
	group.DELETE("/comment/delete", controller.DeleteCommentREST)
	group.POST("/comment/like", controller.LikeCommentREST)
	group.POST("/comment/unlike", controller.UnlikeCommentREST)

	// explore
	exploreGroup := group.Group("/explore")
	exploreGroup.POST("/fill", controller.FillExploreREST)

	// category
	categoryGroup := group.Group("/category")
	categoryGroup.GET("/get", controller.CategoriesREST)
}

func setSearch(server *echo.Echo, controller *search_controller.Controller) {
	group := server.Group("/api/v1/search")

	setMiddlewaresPublic(group)
	group.POST("/user", controller.User)
	group.POST("/category", controller.Category)
	group.POST("/smart", controller.Smart)
}
