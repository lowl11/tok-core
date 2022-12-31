package repositories

import (
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/repositories/auth_repository"
	"tok-core/src/repositories/post_category_repository"
	"tok-core/src/repositories/post_comment_repository"
	"tok-core/src/repositories/post_like_repository"
	"tok-core/src/repositories/post_repository"
	"tok-core/src/repositories/subscription_repository"
	"tok-core/src/repositories/user_ip_repository"
	"tok-core/src/repositories/user_repository"
	"tok-core/src/services/mongo_service"
	"tok-core/src/services/postgres_service"
)

const (
	databaseName = "tok"

	postCommentCollection = "post_comments"
	postLikeCollection    = "post_likes"
)

type ApiRepositories struct {
	Auth         *auth_repository.Repository
	User         *user_repository.Repository
	UserIP       *user_ip_repository.Repository
	Subscription *subscription_repository.Repository

	PostCategory *post_category_repository.Repository
	Post         *post_repository.Repository
	PostComment  *post_comment_repository.Repository
	PostLike     *post_like_repository.Repository
}

func Get(apiEvents *events.ApiEvents) (*ApiRepositories, error) {
	logger := definition.Logger

	// подключение к Postgres
	connectionPostgres, err := postgres_service.NewConnection()
	if err != nil {
		logger.Fatal(err, "Connect to Postgres database error", layers.Database)
	}

	// подключение к MongoDB
	connectionMongo, err := mongo_service.NewConnection(databaseName)
	if err != nil {
		logger.Fatal(err, "Connect to Mongo database error", layers.Mongo)
	}

	return &ApiRepositories{
		Auth:         auth_repository.Create(connectionPostgres, apiEvents),
		User:         user_repository.Create(connectionPostgres, apiEvents),
		UserIP:       user_ip_repository.Create(connectionPostgres, apiEvents),
		Subscription: subscription_repository.Create(connectionPostgres, apiEvents),

		PostCategory: post_category_repository.Create(connectionPostgres, apiEvents),
		Post:         post_repository.Create(connectionPostgres, apiEvents),
		PostComment:  post_comment_repository.Create(connectionMongo, postCommentCollection),
		PostLike:     post_like_repository.Create(connectionMongo, postLikeCollection),
	}, nil
}
