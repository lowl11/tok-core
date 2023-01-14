package models

import (
	"time"
	"tok-core/src/data/interfaces"
)

type NotificationGetInfo struct {
	Username        string             `json:"username"`
	NewActionsCount int                `json:"new_actions_count"`
	Items           []NotificationItem `json:"items"`
}

type NotificationGetCount struct {
	Username        string `json:"username"`
	NewActionsCount int    `json:"new_actions_count"`
}

type NotificationItem struct {
	Status     string                         `json:"status"`
	Username   *UserDynamicGet                `json:"user"`
	ActionKey  string                         `json:"action_key"`
	ActionCode string                         `json:"action_code"`
	ActionBody interfaces.INotificationAction `json:"action_body"`
	CreatedAt  time.Time                      `json:"created_at"`
}
