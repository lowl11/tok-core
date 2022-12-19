package models

type SearchUser struct {
	Query string `json:"query"`
}

type SearchCategory struct {
	Query string `json:"query"`
}

type SearchSmart struct {
	Query string `json:"query"`
}
