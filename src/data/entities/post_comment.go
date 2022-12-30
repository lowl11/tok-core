package entities

import "time"

type PostCommentGet struct {
	PostCode       string    `bson:"post_code"`
	AuthorUsername string    `bson:"author_username"`
	Text           string    `bson:"text"`
	LikesCount     int       `bson:"likes_count"`
	Likers         []string  `bson:"likers"`
	CreatedAt      time.Time `bson:"created_at"`
}

type PostCommentCreate struct {
	PostCode   string `bson:"post_code"`
	PostAuthor string `bson:"post_author"`

	Comments []PostCommentItem `bson:"comments"`
}

type PostCommentItem struct {
	CommentCode   string    `bson:"comment_code"`
	CommentAuthor string    `bson:"comment_author"`
	Text          string    `bson:"text"`
	LikesCount    int       `bson:"likes_count"`
	LikeAuthors   []string  `bson:"like_authors"`
	CreatedAt     time.Time `bson:"created_at"`

	SubComments []PostSubCommentItem `bson:"sub_comments"`
}

type PostSubCommentItem struct {
	CommentCode   string    `bson:"comment_code"`
	CommentAuthor string    `bson:"comment_author"`
	Text          string    `bson:"text"`
	LikesCount    int       `bson:"likes_count"`
	LikeAuthors   []string  `bson:"like_authors"`
	CreatedAt     time.Time `bson:"created_at"`
}
