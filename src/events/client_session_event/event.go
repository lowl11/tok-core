package client_session_event

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Event struct {
	client *redis.Client
}

func Create(address, password string) (*Event, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // default database
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Event{
		client: client,
	}, nil
}
