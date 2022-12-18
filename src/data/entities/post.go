package entities

import "time"

type PostGet struct {
	ID int `db:"id"`

	AuthorUsername string  `db:"author_username"`
	AuthorName     *string `db:"author_name"`
	AuthorAvatar   *string `db:"author_avatar"`

	CategoryName string `db:"category_name"`
	CategoryCode string `db:"category_code"`

	Code    string  `db:"code"`
	Text    string  `db:"text"`
	Picture *string `db:"picture"`

	CreatedAt time.Time `db:"created_at"`
}

type PostCreate struct {
	Code string `db:"code"`

	CategoryCode   string `db:"category_code"`
	AuthorUsername string `db:"author_username"`

	Text    string  `db:"text"`
	Picture *string `db:"picture"`
}
