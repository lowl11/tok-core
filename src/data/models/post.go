package models

import "time"

type PostAdd struct {
	CategoryCode   string  `json:"category_code"`
	CustomCategory *string `json:"custom_category"`

	Text    string `json:"text"`
	Picture *struct {
		Name   string `json:"name"`
		Buffer string `json:"buffer"`
	} `json:"picture"`
}

type PostAddExtended struct {
	Base        *PostAdd
	ImageConfig *ImageConfig
}

type PostGet struct {
	AuthorUsername string  `json:"author_username"`
	AuthorName     *string `json:"author_name"`
	AuthorAvatar   *string `json:"author_avatar"`

	CategoryCode string `json:"category_code"`
	CategoryName string `json:"category_name"`

	Code      string          `json:"code"`
	Text      string          `json:"text"`
	Picture   *PostGetPicture `json:"picture"`
	CreatedAt time.Time       `json:"created_at"`

	LikeCount    int  `json:"like_count"`
	MyLike       bool `json:"my_like"`
	CommentCount int  `json:"comment_count"`
}

type PostGetPicture struct {
	Path   *string `json:"path"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
}
