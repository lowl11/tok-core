package subscription_repository

import (
	"github.com/jmoiron/sqlx"
	"tok-core/src/events"
	"tok-core/src/repositories/repository"
)

type Repository struct {
	repository.Base
	connection *sqlx.DB
}

func Create(connection *sqlx.DB, apiEvents *events.ApiEvents) *Repository {
	return &Repository{
		Base:       repository.CreateBase(apiEvents.Script),
		connection: connection,
	}
}
