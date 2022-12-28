package models

type PostLike struct {
	Author   string `json:"author"`
	PostCode string `json:"post_code"`
}

type PostUnlike struct {
	Author   string `json:"author"`
	PostCode string `json:"post_code"`
}
