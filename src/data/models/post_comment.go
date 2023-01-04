package models

import "time"

type PostCommentItem struct {
	CommentCode   string    `json:"comment_code"`
	CommentAuthor string    `json:"comment_author"`
	CommentText   string    `json:"comment_text"`
	LikesCount    int       `json:"likes_count"`
	LikeAuthors   []string  `json:"like_authors"`
	CreatedAt     time.Time `json:"created_at"`

	SubComments []PostSubCommentItem `json:"subcomments"`
}

type PostSubCommentItem struct {
	CommentCode   string    `json:"comment_code"`
	CommentAuthor string    `json:"comment_author"`
	CommentText   string    `json:"comment_text"`
	LikesCount    int       `json:"likes_count"`
	LikeAuthors   []string  `json:"like_authors"`
	CreatedAt     time.Time `json:"created_at"`
}

type PostCommentGet struct {
	PostCode   string            `json:"post_code"`
	PostAuthor string            `json:"post_author"`
	Comments   []PostCommentItem `json:"comments"`
}

type PostCommentAdd struct {
	PostCode   string `json:"post_code"`
	PostAuthor string `json:"post_author"`

	CommentText string `json:"comment_text"`

	// если комментарий не первый в посте
	ParentCommentCode string `json:"parent_comment_code"`
	FirstComment      bool   `json:"first_comment"`
}

type PostCommentDelete struct {
	CommentCode       string `json:"comment_code"`
	ParentCommentCode string `json:"parent_comment_code"`
	SubComment        bool   `json:"subcomment"`
}

type PostCommentLike struct {
	CommentCode       string `json:"comment_code"`
	ParentCommentCode string `json:"parent_comment_code"`
	SubComment        bool   `json:"subcomment"`
}

type PostCommentUnlike struct {
	CommentCode       string `json:"comment_code"`
	ParentCommentCode string `json:"parent_comment_code"`
	SubComment        bool   `json:"subcomment"`
}
