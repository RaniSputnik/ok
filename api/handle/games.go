package handle

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/game"
)

func OneGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game := model.Game{
			CreatedBy: "RaniSputnik",
			CreatedAt: time.Now().In(time.UTC),
			Black:     "RaniSputnik",
			White:     "Derpzilla",
			Board: model.Board{
				Size: game.BoardSizeSmall,
				Stones: []model.Stone{
					model.Stone{
						Colour: "BLACK",
						X:      3,
						Y:      5,
					},
				},
			},
		}

		json.NewEncoder(w).Encode(&game)
	}
}
