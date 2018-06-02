package handle_test

import (
	"net/http"
	"testing"

	"github.com/RaniSputnik/ok/api/handle"
	"github.com/RaniSputnik/ok/game"
)

func TestIllegalMove(t *testing.T) {
	message := func(msg string) handle.HTTPError {
		return handle.HTTPError{
			Status:  http.StatusBadRequest,
			Type:    "IllegalMove",
			Message: msg,
		}
	}

	testCases := []struct {
		GivenErr error
		Expected handle.HTTPError
	}{
		{
			GivenErr: game.ErrNotYourTurn,
			Expected: message("It is not your turn."),
		},
		{
			GivenErr: game.ErrOutsideBoard,
			Expected: message("Move is outside the board."),
		},
		{
			GivenErr: game.ErrPositionOccupied,
			Expected: message("The position is already occupied."),
		},
		{
			GivenErr: game.ErrSuicidalMove,
			Expected: message("Move would result in suicide."),
		},
		{
			GivenErr: game.ErrViolatesKo,
			Expected: message("Move violates ko."),
		},
	}

	for _, test := range testCases {
		if got := handle.ErrIllegalMove(test.GivenErr); got != test.Expected {
			t.Errorf("Expected: %+v, Got: %+v", test.Expected, got)
		}
	}
}
