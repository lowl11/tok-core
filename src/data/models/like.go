package models

type PostLike struct {
	PostCode   string `json:"post_code"`
	LikeAuthor string `json:"like_author"`
}

type PostUnlike struct {
	PostCode   string `json:"post_code"`
	LikeAuthor string `json:"like_author"`
}
