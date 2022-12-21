package models

type ProfileUpdate struct {
	Name string `json:"name"`
	BIO  string `json:"bio"`
}

type ProfileUpdateContact struct {
	Phone *string `json:"phone"`
	Email *string `json:"email"`
}
