package events

import (
	"tok-core/src/definition"
	"tok-core/src/events/client_session_event"
	"tok-core/src/events/image_event"
	"tok-core/src/events/script_event"
)

type ApiEvents struct {
	Script        *script_event.Event
	ClientSession *client_session_event.Event
	Image         *image_event.Event
}

func Get() (*ApiEvents, error) {
	config := definition.Config

	script, err := script_event.Create()
	if err != nil {
		return nil, err
	}

	clientSession, err := client_session_event.Create(config.Redis.Address, config.Redis.Password)
	if err != nil {
		return nil, err
	}

	image := image_event.Create(config.Image.BasePath)

	return &ApiEvents{
		Script:        script,
		ClientSession: clientSession,
		Image:         image,
	}, nil
}
