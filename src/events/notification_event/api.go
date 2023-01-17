package notification_event

import (
	"github.com/google/uuid"
	"tok-core/src/data/entities"
	"tok-core/src/repositories/notification_repository"
)

func (event *Event) SetRepo(notification *notification_repository.Repository) {
	event.notification = notification
}

func (event *Event) Push(code, username, actionAuthor string, body *entities.NotificationBody) error {
	actionKey := uuid.New().String()

	if err := event.notification.Create(username, actionAuthor, actionKey, code, body); err != nil {
		return err
	}

	return nil
}
