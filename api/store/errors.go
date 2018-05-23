package store

import "errors"

var (
	ErrUsernameTaken = errors.New("username already in use.")
)
