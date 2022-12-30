package models

import "time"

type PostCommentGet struct {
	PostCode       string    `json:"post_code"`
	AuthorUsername string    `json:"author_username"`
	Text           string    `json:"text"`
	LikesCount     int       `json:"likes_count"`
	LikeAuthors    []string  `json:"like_authors"`
	CreatedAt      time.Time `json:"created_at"`
}

type PostCommentAdd struct {
	PostCode   string `json:"post_code"`
	PostAuthor string `json:"post_author"`

	CommentAuthor string `json:"comment_author"`
	CommentText   string `json:"comment_text"`
	FirstComment  bool   `json:"first_comment"`
}
