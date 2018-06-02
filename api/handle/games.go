package handle

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RaniSputnik/ok/api/kontext"
	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/api/store"
	"github.com/RaniSputnik/ok/game"
)

var testGame = model.Game{
	CreatedBy: "Alice",
	CreatedAt: time.Now().In(time.UTC),
	Black:     "Alice",
	White:     "Bob",
	Board: model.Board{
		Size: game.BoardSizeSmall,
		Stones: []model.Stone{
			model.Stone{
				Colour: game.Black.String(),
				X:      3,
				Y:      5,
			},
		},
	},
}

func OneGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&testGame)
	}
}

type RequestVarFunc func(r *http.Request) string

type playInput struct {
	X *int `json:"x"`
	Y *int `json:"y"`
}

func (i playInput) Validate() (err HTTPError, ok bool) {
	if i.X == nil {
		return errMissingParameter("x"), false
	}
	if i.Y == nil {
		return errMissingParameter("y"), false
	}
	return HTTPError{}, true
}

type playResponse struct {
	Colour string `json:"colour"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func Play(getGameID RequestVarFunc, db store.Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		player := kontext.Player(r.Context())
		gameID := getGameID(r)

		var input playInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, errInvalidJSON(err))
			return
		}

		if err, ok := input.Validate(); !ok {
			writeError(w, err)
			return
		}

		g, err := db.GetGameByID(r.Context(), gameID)
		panicIf(err)
		if g == nil {
			writeError(w, errResourceNotFound("game", gameID))
			return
		}

		playerColour := game.None
		if g.Black == player.Username {
			playerColour = game.Black
		} else if g.White == player.Username {
			playerColour = game.White
		} else {
			writeError(w, errNotParticipating())
			return
		}

		moves, err := db.GetGameMoves(r.Context(), gameID)
		panicIf(err)

		// Apply all of the moves from the DB
		// TODO is there a way we can cache this
		// so that we do not have to do this work every
		// time a new move is played and still maintain
		// a robust DB structure?
		m := game.New(g.Board.Size)
		for _, move := range moves {
			panicIf(m.Play(move))
		}

		stone := game.Stone{
			Colour: playerColour,
			Position: game.Position{
				X: *input.X,
				Y: *input.Y,
			}}
		if err := m.Play(stone); err != nil {
			writeError(w, ErrIllegalMove(err))
			return
		}

		panicIf(db.SaveStone(r.Context(), gameID, stone))

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&playResponse{
			Colour: stone.Colour.String(),
			X:      stone.X,
			Y:      stone.Y,
		})
	}
}
