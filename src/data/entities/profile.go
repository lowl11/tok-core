package entities

type ProfileUpdate struct {
	Username string `db:"username"`
	Name     string `db:"name"`
	BIO      string `db:"bio"`
}

type ProfileAvatarUpdate struct {
	Username string `db:"username"`
	Path     string `db:"path"`
}

type ProfileWallpaperUpdate struct {
	Username string `db:"username"`
	Path     string `db:"path"`
}
