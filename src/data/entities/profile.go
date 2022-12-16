package entities

type ProfileUpdate struct {
	Username string `db:"username"`
	Name     string `db:"name"`
	BIO      string `db:"bio"`
}
