package models

type PostLike struct {
	PostCode string `json:"post_code"`
}

type PostUnlike struct {
	PostCode string `json:"post_code"`
}
