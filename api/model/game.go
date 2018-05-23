package model

import "time"

type Game struct {
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	Black     string    `json:"black"`
	White     string    `json:"white"`
	Board     Board     `json:"board"`
}

type Board struct {
	Size   int     `json:"size"`
	Stones []Stone `json:"stones"`
}

type Stone struct {
	Colour string `json:"colour"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
