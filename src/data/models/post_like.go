package models

type PostLikeGet struct {
	LikesCount  int              `json:"likes_count"`
	LikeAuthors []UserDynamicGet `json:"like_authors"`
	Liked       bool             `json:"Liked"`
}

type PostLike struct {
	PostCode   string `json:"post_code"`
	PostAuthor string `json:"post_author"`
}

type PostUnlike struct {
	PostCode string `json:"post_code"`
}
