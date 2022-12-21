package entities

type PostCategoryGet struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

type PostCategoryCreate struct {
	Code string `db:"code"`
	Name string `db:"name"`
}
