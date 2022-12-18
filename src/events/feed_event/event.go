package feed_event

import "github.com/go-redis/redis/v8"

type Event struct {
	client *redis.Client
}

func Create(client *redis.Client) *Event {
	return &Event{
		client: client,
	}
}
