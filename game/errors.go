package game

import "errors"

var (
	// ErrNotYourTurn is returned when the wrong player
	// attempts to make a move.
	ErrNotYourTurn = errors.New("not your turn")

	// ErrOutsideBoard is returned when the players
	// move is outside the board they are playing on.
	ErrOutsideBoard = errors.New("outside board")

	// ErrPositionOccupied is returned when there is
	// already a stone in the move position.
	ErrPositionOccupied = errors.New("position occupied")

	// ErrSuicidalMove is returned when the move would
	// result in the piece being immediately captured.
	// https://senseis.xmp.net/?Suicide
	ErrSuicidalMove = errors.New("suicidal move")

	// ErrViolatesKo is returned when the players move
	// violates the ko rule, repeating board state.
	// https://senseis.xmp.net/?Ko
	ErrViolatesKo = errors.New("violates ko")
)
