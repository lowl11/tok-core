package repositories

import (
	"github.com/jmoiron/sqlx"
	"time"
	"tok-core/src/definition"
	"tok-core/src/events"
	"tok-core/src/repositories/auth_repository"
	"tok-core/src/repositories/subscription_repository"
	"tok-core/src/repositories/user_repository"

	_ "github.com/lib/pq"
)

type ApiRepositories struct {
	Auth         *auth_repository.Repository
	User         *user_repository.Repository
	Subscription *subscription_repository.Repository
}

func Get(apiEvents *events.ApiEvents) (*ApiRepositories, error) {
	config := definition.Config.Database
	logger := definition.Logger

	connection, err := sqlx.Open("postgres", config.Connection)
	if err != nil {
		return nil, err
	}

	connection.SetMaxOpenConns(config.MaxConnections)
	connection.SetMaxIdleConns(config.MaxConnections)
	connection.SetConnMaxIdleTime(time.Duration(config.Lifetime) * time.Minute)

	logger.Info("Ping database...")
	if err = connection.Ping(); err != nil {
		return nil, err
	}
	logger.Info("Ping database done!")

	logger.Info("Initialization database...")
	defer logger.Info("Initialization database done!")

	return &ApiRepositories{
		Auth:         auth_repository.Create(connection, apiEvents),
		User:         user_repository.Create(connection, apiEvents),
		Subscription: subscription_repository.Create(connection, apiEvents),
	}, nil
}
