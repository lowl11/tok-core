package entities

import "time"

type FeedPost struct {
	PostCode     string    `bson:"post_code"`
	PostCategory string    `bson:"post_category"`
	PostAuthor   string    `bson:"post_author"`
	CreatedAt    time.Time `bson:"created_at"`
}

type FeedGet struct {
	Name  string     `bson:"name"`
	Count int        `bson:"count"`
	Posts []FeedPost `bson:"posts"`
}

type FeedCreate struct {
	Name  string     `bson:"name"`
	Count int        `bson:"count"`
	Posts []FeedPost `bson:"posts"`
}
