package models

type PostLikeGet struct {
	LikesCount  int      `json:"likes_count"`
	LikeAuthors []string `json:"like_authors"`
	Liked       bool     `json:"Liked"`
}
