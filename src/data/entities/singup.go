package entities

type Signup struct {
	Username string `db:"username"`
	Password string `db:"password"`
}
