package handle

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RaniSputnik/ok/game"
)

type Game struct {
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	Black     string    `json:"black"`
	White     string    `json:"white"`
	Board     Board     `json:"board"`
}

type Board struct {
	Size   int     `json:"size"`
	Stones []Stone `json:"stones"`
}

type Stone struct {
	Colour string `json:"colour"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func oneGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		game := Game{
			CreatedBy: "RaniSputnik",
			CreatedAt: time.Now().In(time.UTC),
			Black:     "RaniSputnik",
			White:     "Derpzilla",
			Board: Board{
				Size: game.BoardSizeSmall,
				Stones: []Stone{
					Stone{
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
