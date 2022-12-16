package auth_repository

import (
	"github.com/jmoiron/sqlx"
	"tok-core/src/events"
	"tok-core/src/events/script_event"
	"tok-core/src/repositories/repository"
)

type Repository struct {
	repository.Base
	connection *sqlx.DB
	script     *script_event.Event
}

func Create(connection *sqlx.DB, apiEvents *events.ApiEvents) *Repository {
	return &Repository{
		connection: connection,
		script:     apiEvents.Script,
	}
}
