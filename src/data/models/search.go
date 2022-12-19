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

type SearchItemGet struct {
	User     *SearchUserGet     `json:"user"`
	Category *SearchCategoryGet `json:"category"`
}

type SearchCategoryGet struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type SearchUserGet struct {
	Username string  `json:"username"`
	Name     *string `json:"name"`
	Avatar   *string `json:"avatar"`

	// подписан ли авторизованный
	Subscribed bool `json:"subscribed"`
}
