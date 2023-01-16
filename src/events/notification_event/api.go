package notification_event

import (
	"github.com/google/uuid"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/repositories/notification_repository"
)

func (event *Event) SetRepo(notification *notification_repository.Repository) {
	event.notification = notification
}

func (event *Event) Push(code, username string, action any) error {
	actionKey := uuid.New().String()

	if err := event.notification.AddItemExist(username, &entities.NotificationAction{
		Username:   username,
		Status:     "new",
		ActionKey:  actionKey,
		ActionCode: code,
		ActionBody: action,
		CreatedAt:  time.Now(),
	}); err != nil {
		return err
	}

	return nil
}
