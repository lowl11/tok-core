package action_helper

import (
	"tok-core/src/data/interfaces"
	"tok-core/src/data/notification_actions"
)

func Get(actionCode string) interfaces.INotificationAction {
	switch actionCode {
	case PostLike:
		return &notification_actions.PostLike{}
	}

	return nil
}
