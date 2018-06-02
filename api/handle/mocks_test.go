package handle_test

import (
	"context"

	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/game"
)

type mockPlayerStore struct {
	Func struct {
		SavePlayer struct {
			Called struct {
				With struct {
					Ctx    context.Context
					Player *model.Player
				}
				Times int
			}
			Returns struct {
				Err error
			}
		}
	}
}

func (m *mockPlayerStore) SavePlayer(ctx context.Context, player *model.Player) error {
	m.Func.SavePlayer.Called.Times++
	m.Func.SavePlayer.Called.With.Ctx = ctx
	m.Func.SavePlayer.Called.With.Player = player
	return m.Func.SavePlayer.Returns.Err
}

type mockGameStore struct {
	Func struct {
		GetGameByID struct {
			Called struct {
				With struct {
					Ctx context.Context
					ID  string
				}
				Times int
			}
			Returns struct {
				Game *model.Game
				Err  error
			}
		}
		GetGameMoves struct {
			Called struct {
				With struct {
					Ctx    context.Context
					GameID string
				}
				Times int
			}
			Returns struct {
				Moves []game.Move
				Err   error
			}
		}
		SaveStone struct {
			Called struct {
				With struct {
					Ctx    context.Context
					GameID string
					Stone  game.Stone
				}
				Times int
			}
			Returns struct {
				Err error
			}
		}
	}
}

func (m *mockGameStore) GetGameByID(ctx context.Context, gameID string) (*model.Game, error) {
	m.Func.GetGameByID.Called.Times++
	m.Func.GetGameByID.Called.With.Ctx = ctx
	m.Func.GetGameByID.Called.With.ID = gameID
	return m.Func.GetGameByID.Returns.Game, m.Func.GetGameByID.Returns.Err
}

func (m *mockGameStore) GetGameMoves(ctx context.Context, gameID string) ([]game.Move, error) {
	m.Func.GetGameMoves.Called.Times++
	m.Func.GetGameMoves.Called.With.Ctx = ctx
	m.Func.GetGameMoves.Called.With.GameID = gameID
	return m.Func.GetGameMoves.Returns.Moves, m.Func.GetGameMoves.Returns.Err
}

func (m *mockGameStore) SaveStone(ctx context.Context, gameID string, stone game.Stone) error {
	m.Func.SaveStone.Called.Times++
	m.Func.SaveStone.Called.With.Ctx = ctx
	m.Func.SaveStone.Called.With.GameID = gameID
	m.Func.SaveStone.Called.With.Stone = stone
	return m.Func.SaveStone.Returns.Err
}
