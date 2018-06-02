package kontext

import (
	"context"

	"github.com/RaniSputnik/ok/api/model"
)

func WithPlayer(parent context.Context, player *model.Player) context.Context {
	return context.WithValue(parent, model.Player{}, player)
}

func Player(from context.Context) *model.Player {
	val, _ := from.Value(model.Player{}).(*model.Player)
	return val
}
