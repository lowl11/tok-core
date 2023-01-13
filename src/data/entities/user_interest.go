package entities

type UserInterestCategory struct {
	CategoryCode string `bson:"category_code"`
	Interest     int    `bson:"interest"`
}

type UserInterestGet struct {
	Username   string                 `bson:"username"`
	Categories []UserInterestCategory `bson:"categories"`
}

type UserInterestCreate struct {
	Username   string                 `bson:"username"`
	Categories []UserInterestCategory `bson:"categories"`
}
