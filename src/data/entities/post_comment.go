package entities

import "time"

type PostCommentItem struct {
	CommentCode   string    `bson:"comment_code"`
	CommentAuthor string    `bson:"comment_author"`
	CommentText   string    `bson:"comment_text"`
	LikesCount    int       `bson:"likes_count"`
	LikeAuthors   []string  `bson:"like_authors"`
	CreatedAt     time.Time `bson:"created_at"`

	SubComments []PostSubCommentItem `bson:"subcomments"`
}

type PostSubCommentItem struct {
	CommentCode   string    `bson:"comment_code"`
	CommentAuthor string    `bson:"comment_author"`
	CommentText   string    `bson:"comment_text"`
	LikesCount    int       `bson:"likes_count"`
	LikeAuthors   []string  `bson:"like_authors"`
	CreatedAt     time.Time `bson:"created_at"`
}

type PostCommentCreate struct {
	PostCode   string `bson:"post_code"`
	PostAuthor string `bson:"post_author"`

	Comments []PostCommentItem `bson:"comments"`
}

type PostCommentAppendComment struct {
	CommentCode   string    `bson:"comment_code"`
	CommentAuthor string    `bson:"comment_author"`
	CommentText   string    `bson:"comment_text"`
	LikesCount    int       `bson:"likes_count"`
	LikeAuthors   []string  `bson:"like_authors"`
	CreatedAt     time.Time `bson:"created_at"`

	SubComments []PostSubCommentItem `bson:"subcomments"`
}

type PostCommentAppendSubComment struct {
	CommentCode   string    `bson:"comment_code"`
	CommentAuthor string    `bson:"comment_author"`
	CommentText   string    `bson:"comment_text"`
	LikesCount    int       `bson:"likes_count"`
	LikeAuthors   []string  `bson:"like_authors"`
	CreatedAt     time.Time `bson:"created_at"`
}

type PostCommentGet struct {
	PostCode       string    `bson:"post_code"`
	AuthorUsername string    `bson:"author_username"`
	Text           string    `bson:"text"`
	LikesCount     int       `bson:"likes_count"`
	Likers         []string  `bson:"likers"`
	CreatedAt      time.Time `bson:"created_at"`
}
