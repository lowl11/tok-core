package entities

type ClientSession struct {
	ID int `json:"id"`

	// main
	Username string `json:"username"`
	Password string `json:"password"`

	// additional
	Name *string `json:"name"`
	BIO  *string `json:"bio"`

	// images
	Avatar    *string `json:"avatar"`
	Wallpaper *string `json:"wallpaper"`

	// subscriptions
	Subscriptions ClientSessionSubscribes `json:"subscriptions"`
}

type ClientSessionSubscribes struct {
	SubscriberCount   int `json:"subscriber_count"`
	SubscriptionCount int `json:"subscription_count"`

	Subscribers   []string `json:"subscribers"`
	Subscriptions []string `json:"subscriptions"`
}
