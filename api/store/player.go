package store

import (
	"context"

	"github.com/RaniSputnik/ok/api/model"
)

type Player interface {
	SavePlayer(ctx context.Context, input *model.Player) error
}
