package kontext_test

import (
	"context"
	"testing"

	"github.com/RaniSputnik/ok/api/kontext"
	"github.com/RaniSputnik/ok/api/model"
	"github.com/stretchr/testify/assert"
)

func TestPlayer(t *testing.T) {
	givenPlayer := &model.Player{
		Username: "Clive",
	}

	t.Run("CanBeRetrieved", func(t *testing.T) {
		ctx := context.Background()
		ctxWithPlayer := kontext.WithPlayer(ctx, givenPlayer)
		gotPlayer := kontext.Player(ctxWithPlayer)

		assert.NotEqual(t, ctx, ctxWithPlayer)
		assert.Equal(t, givenPlayer, gotPlayer)
	})

	t.Run("IsNilWhenUnset", func(t *testing.T) {
		ctx := context.Background()
		gotPlayer := kontext.Player(ctx)

		assert.Nil(t, gotPlayer)
	})
}
