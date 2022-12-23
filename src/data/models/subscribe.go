package models

type ProfileSubscribe struct {
	Username string `json:"username"`
}

type ProfileUnsubscribe struct {
	Username string `json:"username"`
}
