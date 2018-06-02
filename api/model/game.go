package model

import (
	"time"
)

type Game struct {
	ID        string    `json:"id"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	Black     string    `json:"black"`
	White     string    `json:"white"`
	Board     Board     `json:"board"`

	Moves []Move `json:"moves"`
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

type MoveKind string

const (
	MoveKindState = MoveKind("SetState")
	MoveKindStone = MoveKind("PlayStone")
	MoveKindSkip  = MoveKind("SkipTurn")
)

// TODO possibly break into three types
// or use interface

type Move struct {
	// Common properties
	Kind    MoveKind `json:"kind"`
	Colour  string   `json:"colour"`
	Message string   `json:"message"`
	// SetState
	State []byte `json:"state,omitempty"`
	// PlayStone
	X *int `json:"x,omitempty"`
	Y *int `json:"y,omitempty"`
}
