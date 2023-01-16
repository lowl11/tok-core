package models

import (
	"time"
)

type NotificationGetInfo struct {
	Username        string               `json:"username"`
	NewActionsCount int                  `json:"new_actions_count"`
	Actions         []NotificationAction `json:"actions"`
}

type NotificationGetCount struct {
	Username        string `json:"username"`
	NewActionsCount int    `json:"new_actions_count"`
}

type NotificationRead struct {
	ActionKeys []string `json:"action_keys"`
}

type NotificationAction struct {
	Status     string          `json:"status"`
	User       *UserDynamicGet `json:"user"`
	ActionKey  string          `json:"action_key"`
	ActionCode string          `json:"action_code"`
	ActionBody any             `json:"action_body"`
	CreatedAt  time.Time       `json:"created_at"`
}
