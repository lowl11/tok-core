package entities

import "time"

type NotificationAction struct {
	Username   string    `bson:"username"`
	Status     string    `bson:"status"`
	ActionKey  string    `bson:"action_key"`
	ActionCode string    `bson:"action_code"`
	ActionBody any       `bson:"action_body"`
	CreatedAt  time.Time `bson:"created_at"`
}

type NotificationGetInfo struct {
	Username        string               `bson:"username"`
	NewActionsCount int                  `bson:"new_actions_count"`
	Actions         []NotificationAction `bson:"actions"`
}

type NotificationGetCount struct {
	Username        string `bson:"username"`
	NewActionsCount int    `bson:"new_actions_count"`
}

type NotificationCreate struct {
	Username        string               `bson:"username"`
	NewActionsCount int                  `bson:"new_actions_count"`
	Actions         []NotificationAction `bson:"actions"`
}
