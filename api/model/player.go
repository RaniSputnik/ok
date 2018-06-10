package model

type Player struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}
