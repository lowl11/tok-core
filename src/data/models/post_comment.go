package models

import "time"

type PostCommentItem struct {
	CommentCode   string          `json:"comment_code"`
	CommentAuthor *UserDynamicGet `json:"comment_author"`
	CommentText   string          `json:"comment_text"`
	LikesCount    int             `json:"likes_count"`
	MyLike        bool            `json:"my_like"`
	CreatedAt     time.Time       `json:"created_at"`

	SubCommentsSize int `json:"subcomments_size"`
}

type PostSubCommentItem struct {
	CommentCode   string          `json:"comment_code"`
	CommentAuthor *UserDynamicGet `json:"comment_author"`
	CommentText   string          `json:"comment_text"`
	LikesCount    int             `json:"likes_count"`
	MyLike        bool            `json:"my_like"`
	CreatedAt     time.Time       `json:"created_at"`
}

type PostCommentGet struct {
	PostCode   string            `json:"post_code"`
	PostAuthor string            `json:"post_author"`
	Comments   []PostCommentItem `json:"comments"`
}

type PostSubcommentGet struct {
	PostCode    string               `json:"post_code"`
	PostAuthor  string               `json:"post_author"`
	Subcomments []PostSubCommentItem `json:"subcomments"`
}

type PostCommentAdd struct {
	PostCode   string `json:"post_code"`
	PostAuthor string `json:"post_author"`

	CommentText string `json:"comment_text"`

	// Если комментарий не первый в посте
	ParentCommentCode string `json:"parent_comment_code"`
}

type PostCommentDelete struct {
	CommentCode       string `json:"comment_code"`
	ParentCommentCode string `json:"parent_comment_code"`
	SubComment        bool   `json:"subcomment"`
}

type PostCommentLike struct {
	PostCode          string `json:"post_code"`
	CommentCode       string `json:"comment_code"`
	ParentCommentCode string `json:"parent_comment_code"`
	SubComment        bool   `json:"subcomment"`
}

type PostCommentUnlike struct {
	PostCode          string `json:"post_code"`
	CommentCode       string `json:"comment_code"`
	ParentCommentCode string `json:"parent_comment_code"`
	SubComment        bool   `json:"subcomment"`
}
