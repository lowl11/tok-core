package repositories

import (
	"github.com/jmoiron/sqlx"
	"tok-core/src/definition"
	"tok-core/src/events"
	"time"

	_ "github.com/lib/pq"
)

type ApiRepositories struct {
	//
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
		//
	}, nil
}