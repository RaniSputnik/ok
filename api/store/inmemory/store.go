package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/api/store"
	"github.com/RaniSputnik/ok/game"
	"github.com/rs/xid"
)

func New() store.Store {
	return &inMemory{
		players: map[string]*model.Player{},
		games:   map[string]*model.Game{},
		moves:   map[string][]game.Move{},
	}
}

type inMemory struct {
	players map[string]*model.Player
	games   map[string]*model.Game
	moves   map[string][]game.Move

	sync.Mutex
}

func (s *inMemory) GetPlayer(ctx context.Context, username string) (*model.Player, error) {
	s.Lock()
	player := s.players[username]
	s.Unlock()
	return player, nil
}

func (s *inMemory) SavePlayer(ctx context.Context, input *model.Player) error {
	s.Lock()
	if existing := s.players[input.Username]; existing != nil {
		return store.ErrUsernameTaken
	}
	s.players[input.Username] = input
	s.Unlock()
	return nil
}

func (s *inMemory) GetGameByID(ctx context.Context, id string) (*model.Game, error) {
	s.Lock()
	g := s.games[id]
	s.Unlock()
	return g, nil
}

func (s *inMemory) GetGameMoves(ctx context.Context, gameID string) ([]game.Move, error) {
	s.Lock()
	defer s.Unlock()
	if g := s.games[gameID]; g == nil {
		return nil, errors.New("game not found")
	}
	if _, ok := s.moves[gameID]; !ok {
		s.moves[gameID] = []game.Move{}
	}
	return s.moves[gameID], nil
}

func (s *inMemory) SaveGame(ctx context.Context, game *model.Game) error {
	s.Lock()
	defer s.Unlock()
	if game.ID == "" {
		game.ID = xid.New().String()
	}
	s.games[game.ID] = game
	return nil
}

func (s *inMemory) SaveStone(ctx context.Context, gameID string, move game.Stone) error {
	s.Lock()
	defer s.Unlock()
	if g := s.games[gameID]; g == nil {
		return errors.New("game not found")
	}
	if _, ok := s.moves[gameID]; !ok {
		s.moves[gameID] = []game.Move{}
	}
	s.moves[gameID] = append(s.moves[gameID], move)
	return nil
}
