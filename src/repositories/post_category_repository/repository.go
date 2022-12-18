package post_category_repository

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
		connection: connection,
		Base:       repository.CreateBase(apiEvents.Script),
	}
}
