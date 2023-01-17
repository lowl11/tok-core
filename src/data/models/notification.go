package models

import (
	"time"
)

type NotificationGet struct {
	Status       string            `json:"status"`
	User         *UserDynamicGet   `json:"user"`
	ActionAuthor string            `json:"action_author"`
	ActionKey    string            `json:"action_key"`
	ActionCode   string            `json:"action_code"`
	ActionBody   *NotificationBody `json:"action_body"`
	CreatedAt    time.Time         `json:"created_at"`
}

type NotificationRead struct {
	ActionKeys []string `json:"action_keys"`
}

type NotificationPostGet struct {
	Code  string  `json:"code"`
	Image *string `json:"image"`
	Text  string  `json:"text"`
}

type NotificationBody struct {
	Post *NotificationPostGet `json:"post"`
}
