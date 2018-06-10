package store

import "errors"

type Store interface {
	Player
	Game
}

var (
	ErrUsernameTaken = errors.New("username already in use.")
)
