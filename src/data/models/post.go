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

type PostElasticAdd struct {
	Code         string          `json:"code"`
	Text         string          `json:"text"`
	Category     string          `json:"category"`
	CategoryName string          `json:"category_name"`
	Picture      *PostGetPicture `json:"picture"`
	Author       string          `json:"author"`
	CreatedAt    time.Time       `json:"created_at"`

	Keys []string `json:"keys"`
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
}

type PostElasticGet struct {
	Code         string          `json:"code"`
	Text         string          `json:"text"`
	Category     string          `json:"category"`
	CategoryName string          `json:"category_name"`
	Picture      *PostGetPicture `json:"picture"`
	Author       string          `json:"author"`

	Keys []string `json:"keys"`

	CreatedAt time.Time `json:"created_at"`
}

type PostGetPicture struct {
	Path   *string `json:"path"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
}
