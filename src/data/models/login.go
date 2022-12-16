package models

type LoginByCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginByToken struct {
	Token string `json:"token"`
}
