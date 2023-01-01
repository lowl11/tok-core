package entities

type PostLikeGet struct {
	PostCode    string   `bson:"post_code"`
	PostAuthor  string   `bson:"post_author"`
	LikesCount  int      `bson:"likes_count"`
	LikeAuthors []string `bson:"like_authors"`
}

type PostLikeGetList struct {
	PostCode    string   `bson:"post_code"`
	LikesCount  int      `bson:"likes_count"`
	LikeAuthors []string `bson:"like_authors"`
}
