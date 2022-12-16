package models

type ClientSessionCreate struct {
	// main
	Username string `json:"username"`

	// additional
	Name *string `json:"name"`
	BIO  string  `json:"bio"`

	// images
	Avatar    *string `json:"avatar"`
	Wallpaper *string `json:"wallpaper"`
}

type ClientSessionGet struct {
	Username string  `json:"username"`
	Name     *string `json:"name"`
}
