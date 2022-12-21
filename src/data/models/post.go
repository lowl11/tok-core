package models

import "time"

type PostAdd struct {
	CategoryCode string `json:"category_code"`

	Text    string `json:"text"`
	Picture *struct {
		Name   string `json:"name"`
		Buffer string `json:"buffer"`
	} `json:"picture"`
}

type PostGet struct {
	AuthorUsername string  `json:"author_username"`
	AuthorName     *string `json:"author_name"`
	AuthorAvatar   *string `json:"author_avatar"`

	CategoryCode string `json:"category_code"`
	CategoryName string `json:"category_name"`

	Code      string    `json:"code"`
	Text      string    `json:"text"`
	Picture   *string   `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}
