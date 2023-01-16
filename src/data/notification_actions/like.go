package notification_actions

type PostLike struct {
	PostCode string `json:"post_code" bson:"post_code"`
}

func (action *PostLike) Body() any {
	return action
}
