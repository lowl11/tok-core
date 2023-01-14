package models

type UserSubscriber struct {
	Username string  `json:"username"`
	Name     *string `json:"name"`
	Avatar   *string `json:"avatar"`

	// Подписан ли авторизованный
	Subscribed bool `json:"subscribed"`
}

type UserSubscriptions struct {
	Username string  `json:"username"`
	Name     *string `json:"name"`
	Avatar   *string `json:"avatar"`

	// Подписан ли авторизованный
	Subscribed bool `json:"subscribed"`
}

type UserDynamicGet struct {
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
	Name     *string `json:"name"`

	// Подписан ли авторизованный
	Subscribed *bool `json:"subscribed,omitempty"`
}
