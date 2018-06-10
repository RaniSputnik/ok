package store

import (
	"context"

	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/game"
)

type Game interface {
	GetGameByID(ctx context.Context, gameID string) (*model.Game, error)
	GetGameMoves(ctx context.Context, gameID string) ([]game.Move, error)
	SaveGame(ctx context.Context, game *model.Game) error
	SaveStone(ctx context.Context, gameID string, move game.Stone) error
}
