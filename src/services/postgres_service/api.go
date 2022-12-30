package postgres_service

import (
	"github.com/jmoiron/sqlx"
	"time"
	"tok-core/src/definition"

	_ "github.com/lib/pq"
)

func NewConnection() (*sqlx.DB, error) {
	config := definition.Config.Database
	logger := definition.Logger

	// подключение к Postgres
	connection, err := sqlx.Open("postgres", config.Connection)
	if err != nil {
		return nil, err
	}

	connection.SetMaxOpenConns(config.MaxConnections)
	connection.SetMaxIdleConns(config.MaxConnections)
	connection.SetConnMaxIdleTime(time.Duration(config.Lifetime) * time.Minute)

	logger.Info("Ping Postgres database...")
	if err = connection.Ping(); err != nil {
		return nil, err
	}
	logger.Info("Ping Postgres database done!")

	return connection, nil
}
