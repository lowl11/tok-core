package entities

type UserGet struct {
	ID int `db:"id"`

	// main
	Username string `db:"username"`
	Password string `db:"password"`

	// info
	Name *string `db:"name"`
	BIO  *string `db:"bio"`

	// images
	Avatar    *string `db:"avatar"`
	Wallpaper *string `db:"wallpaper"`
}
