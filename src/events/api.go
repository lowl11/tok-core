package events

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
	"tok-core/src/definition"
	"tok-core/src/events/client_session_event"
	"tok-core/src/events/image_event"
	"tok-core/src/events/notification_event"
	"tok-core/src/events/post_category_event"
	"tok-core/src/events/script_event"
)

type ApiEvents struct {
	Script *script_event.Event
	Image  *image_event.Event

	ClientSession *client_session_event.Event

	PostCategory *post_category_event.Event
	Notification *notification_event.Event
}

func Get() (*ApiEvents, error) {
	config := definition.Config

	// redis connection
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       0, // default database
	})

	// ping redis
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	script, err := script_event.Create()
	if err != nil {
		return nil, err
	}

	image := image_event.Create(config.Image.BasePath)

	clientSession := client_session_event.Create(client)

	postCategory := post_category_event.Create()
	notification := notification_event.Create()

	return &ApiEvents{
		Script: script,
		Image:  image,

		ClientSession: clientSession,

		PostCategory: postCategory,
		Notification: notification,
	}, nil
}
