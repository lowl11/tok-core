package models

type PostAdd struct {
	Text         string  `json:"text"`
	Picture      *string `json:"picture"`
	CategoryCode string  `json:"category_code"`
}
