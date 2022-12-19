package client_session_event

import (
	"github.com/go-redis/redis/v8"
	"tok-core/src/events/redis_event"
)

type Event struct {
	redis_event.Base
	client *redis.Client
}

func Create(client *redis.Client) *Event {
	return &Event{
		client: client,
	}
}
