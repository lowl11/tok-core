package entities

type PostCategoryGet struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

type PostCategoryCreate struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

type PostCategoryCountCreate struct {
	CategoryCode string `db:"category_code"`
	Count        int    `db:"count"`
}

type PostCategoryCountGet struct {
	CategoryCode string `db:"category_code"`
	Count        int    `db:"count"`
}
