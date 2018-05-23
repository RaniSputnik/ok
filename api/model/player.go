package model

type Player struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
	Verified  bool   `json:"verified"`
}
