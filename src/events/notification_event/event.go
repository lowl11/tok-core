package notification_event

import "tok-core/src/repositories/notification_repository"

type Event struct {
	notification *notification_repository.Repository
}

func Create() *Event {
	return &Event{}
}
