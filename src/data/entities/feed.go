package entities

import "time"

type FeedPostPicture struct {
	Path   string `bson:"path"`
	Width  int    `bson:"width"`
	Height int    `bson:"height"`
}

type FeedPost struct {
	PostCode     string `bson:"post_code"`
	PostCategory string `bson:"post_category"`

	AuthorUsername string  `bson:"author_username"`
	AuthorName     *string `bson:"author_name"`
	AuthorAvatar   *string `bson:"author_avatar"`

	PostText string           `bson:"post_text"`
	Picture  *FeedPostPicture `bson:"picture"`

	CreatedAt time.Time `bson:"created_at"`
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
