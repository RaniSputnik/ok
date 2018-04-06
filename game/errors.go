package game

import "errors"

var (
	// ErrNotYourTurn is returned when the wrong player
	// attempts to make a move.
	ErrNotYourTurn = errors.New("not your turn")
)
