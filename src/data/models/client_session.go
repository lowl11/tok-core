package models

import "tok-core/src/data/entities"

type ClientSessionCreate struct {
	// main
	Username string `json:"username"`

	// additional
	Name *string `json:"name"`
	BIO  *string `json:"bio"`

	// images
	Avatar    *string `json:"avatar"`
	Wallpaper *string `json:"wallpaper"`
}

type ClientSessionGet struct {
	Token    string `json:"token"`
	Username string `json:"username"`

	Name *string `json:"name"`
	BIO  *string `json:"bio"`

	Avatar    *string `json:"avatar"`
	Wallpaper *string `json:"wallpaper"`

	Subscriptions entities.ClientSessionSubscribes `json:"subscriptions"`
}

type UserInfoGet struct {
	MySubscription bool `json:"my_subscription"`

	Name *string `json:"name"`
	BIO  *string `json:"bio"`

	Avatar    *string `json:"avatar"`
	Wallpaper *string `json:"wallpaper"`

	Subscriptions entities.ClientSessionSubscribes `json:"subscriptions"`
}
