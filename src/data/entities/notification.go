package entities

import "time"

type NotificationGet struct {
	Username   string            `bson:"username"`
	Status     string            `bson:"status"`
	ActionKey  string            `bson:"action_key"`
	ActionCode string            `bson:"action_code"`
	ActionBody *NotificationBody `bson:"action_body"`
	CreatedAt  time.Time         `bson:"created_at"`
}

type NotificationCreate struct {
	Username   string            `bson:"username"`
	Status     string            `bson:"status"`
	ActionKey  string            `bson:"action_key"`
	ActionCode string            `bson:"action_code"`
	ActionBody *NotificationBody `bson:"action_body"`
	CreatedAt  time.Time         `bson:"created_at"`
}

type NotificationBody struct {
	PostCode *string `bson:"post_code"`
}
