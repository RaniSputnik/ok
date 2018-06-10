package model

type RegisterParams struct {
	Username string `json:"username"`
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string  `json:"token"`
	Player *Player `json:"player"`
}
