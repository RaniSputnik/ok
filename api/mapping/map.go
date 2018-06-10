package mapping

import (
	"fmt"

	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/game"
)

func FromModelGame(g model.Game) *game.Match {
	return game.New(g.Board.Size)
}

func ToModelStones(m *game.Match) []model.Stone {
	res := []model.Stone{}
	board := m.Board()
	for i, stone := range board.Stones {
		if stone != game.None {
			x, y := i%board.Size, i/board.Size
			res = append(res, model.Stone{
				Colour: stone.String(),
				X:      x,
				Y:      y,
			})
		}
	}
	return res
}

func ToModelMoves(moves []game.Move) []model.Move {
	res := make([]model.Move, len(moves))
	for i, mv := range moves {
		res[i] = ToModelMove(mv)
	}
	return res
}

func ToModelMove(move game.Move) model.Move {
	switch v := move.(type) {
	case game.Stone:
		return model.Move{
			Kind:    model.MoveKindStone,
			Colour:  v.Colour.String(),
			Message: fmt.Sprintf("%s stone played at %s", v.Colour, Position(v.X, v.Y)),
			X:       &v.X,
			Y:       &v.Y,
		}
	}

	panic(fmt.Sprintf("unrecognised move type: %+v", move))
}

// Position accepts an x, y position on the board
// and returns the human readable, descriptive location
// of the position.
//
// eg. 2,5 becomes C6)
func Position(x, y int) string {
	const xlabels = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return fmt.Sprintf("%s%d", string(xlabels[x]), y+1)
}
