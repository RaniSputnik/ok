package store

import (
	"context"

	"github.com/RaniSputnik/ok/api/model"
)

type Store interface {
	Player
}

func NewInMemory() Store {
	return &inMemory{
		players: map[string]*model.Player{},
	}
}

type inMemory struct {
	players map[string]*model.Player
}

func (s *inMemory) SavePlayer(ctx context.Context, input *model.Player) error {
	if existing := s.players[input.Username]; existing != nil {
		return ErrUsernameTaken
	}
	s.players[input.Username] = input
	return nil
}
