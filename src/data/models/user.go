package models

type UserSubscriber struct {
	Username string  `json:"username"`
	Name     *string `json:"name"`
	Avatar   *string `json:"avatar"`

	// подписан ли авторизованный
	Subscribed bool `json:"subscribed"`
}

type UserSubscriptions struct {
	Username string  `json:"username"`
	Name     *string `json:"name"`
	Avatar   *string `json:"avatar"`

	// подписан ли авторизованный
	Subscribed bool `json:"subscribed"`
}
