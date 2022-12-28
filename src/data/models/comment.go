package models

type PostCommentAdd struct {
	Author   string `json:"author"`
	Text     string `json:"text"`
	PostCode string `json:"post_code"`
}
