package events

import (
	"tok-core/src/events/script_event"
)

type ApiEvents struct {
	Script *script_event.Event
}

func Get() (*ApiEvents, error) {
	script, err := script_event.Create()
if err != nil {
	return nil, err
}

	return &ApiEvents{
		Script: script,
	}, nil
}