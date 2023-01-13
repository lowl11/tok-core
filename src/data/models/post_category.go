package models

type PostCategoryGet struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type PostCategoryAdd struct {
	Name string `json:"name"`
}
