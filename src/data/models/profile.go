package models

type ProfileUpdate struct {
	Name string `json:"name"`
	BIO  string `json:"bio"`
}
