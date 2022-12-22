package models

type LoginByCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginByToken struct {
	Token string `json:"token"`
}
