package store

import (
	"context"

	"github.com/RaniSputnik/ok/api/model"
)

type Player interface {
	GetPlayer(ctx context.Context, username string) (*model.Player, error)
	SavePlayer(ctx context.Context, input *model.Player) error
}
