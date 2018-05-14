package game_test

import (
	"testing"

	"github.com/RaniSputnik/ko2/game"
)

func TestNewMatchSetsSize(t *testing.T) {
	testCases := []int{game.BoardSizeTiny, game.BoardSizeSmall, game.BoardSizeNormal}

	for _, size := range testCases {
		t.Logf("Testing size: '%d'", size)

		m := game.New(size)
		if gotSize := m.Board().Size; gotSize != size {
			t.Errorf("Expected size: '%d', got: '%d'", size, gotSize)
		}

		if gotStonesLen := len(m.Board().Stones); gotStonesLen != size*size {
			t.Errorf("Expected stones length: '%d', got: '%d'", size*size, gotStonesLen)
		}
	}
}

func TestNewMatchCreatesStones(t *testing.T) {
	stones := []game.Stone{
		black(0, 0),
		black(7, 1),
		white(1, 0),
		white(5, 4),
		white(3, 7),
	}

	m := game.New(game.BoardSizeTiny, stones...)
	got := m.Board()

	for _, stone := range stones {
		if gotCol := got.At(stone.X, stone.Y); gotCol != stone.Colour {
			t.Errorf("Expected: '%s' at {%d,%d}, got: '%s'",
				stone.Colour, stone.X, stone.Y, gotCol)
		}
	}
}

func TestPlayAddsAStone(t *testing.T) {
	m := game.New(game.BoardSizeTiny)
	playedStone := black(0, 0)
	err := m.Play(playedStone)

	if err != nil {
		t.Errorf("Expected err: '<nil>', got: '%v'", err)
	}

	gotStone := m.Board().Stones[0]
	if gotStone != playedStone.Colour {
		t.Errorf("Expected: '%s', got: '%s'", playedStone.Colour, gotStone)
	}
}

func TestBlackPlaysFirst(t *testing.T) {
	m := game.New(game.BoardSizeTiny)
	if got := m.Next(); got != game.Black {
		t.Errorf("Expected: '%s', got: '%s'", game.Black, got)
	}
}

func TestPlayFailsWhenNotYourTurn(t *testing.T) {
	m := game.New(game.BoardSizeTiny)
	err := m.Play(white(0, 0))
	if err != game.ErrNotYourTurn {
		t.Errorf("Expected: '%v', got: '%v'", game.ErrNotYourTurn, err)
	}
}

func TestPlayFailsWhenMoveIsOutsideBoard(t *testing.T) {
	testCases := []struct {
		BoardSize    int
		PlayX, PlayY int
	}{
		{BoardSize: 9, PlayX: -1, PlayY: -1},
		{BoardSize: 9, PlayX: -1, PlayY: 0},
		{BoardSize: 2, PlayX: 5, PlayY: -1},
		{BoardSize: 2, PlayX: 5, PlayY: 0},
		{BoardSize: 5, PlayX: 5, PlayY: 0},
		{BoardSize: 5, PlayX: 3, PlayY: 5},
		{BoardSize: 9, PlayX: 0, PlayY: 100},
	}

	expected := game.ErrOutsideBoard
	for _, test := range testCases {
		m := game.New(test.BoardSize)
		err := m.Play(black(test.PlayX, test.PlayY))
		if err != expected {
			t.Errorf("Expected: '%v', Got: '%v'. %+v", expected, err, test)
		}
	}
}

func TestPlayChangesNext(t *testing.T) {
	m := game.New(game.BoardSizeTiny)
	m.Play(black(0, 0))
	if got := m.Next(); got != game.White {
		t.Errorf("Expected: '%s' to go second, got: '%s'", game.White, got)
	}
	m.Play(white(1, 0))
	if got := m.Next(); got != game.Black {
		t.Errorf("Expected: '%s' to go third, got: '%s'", game.Black, got)
	}
}

func TestSurroundedStonesAreCaptured(t *testing.T) {
	stones := []game.Stone{black(2, 2), black(1, 3), black(3, 3), white(2, 3)}
	m := game.New(game.BoardSizeTiny, stones...)
	m.Play(black(2, 4))

	if got := m.Board().At(2, 3); got != game.None {
		t.Errorf("Expected white stone at position {2,3} to be captured, instead got: '%s'", got)
	}
}

func TestCorneredStonesAreCaptured(t *testing.T) {
	stones := []game.Stone{black(0, 1), white(0, 0)}
	m := game.New(game.BoardSizeTiny, stones...)
	m.Play(black(1, 0))

	if got := m.Board().At(0, 0); got != game.None {
		t.Errorf("Expected white stone at position {0,0} to be captured, instead got: '%s'", got)
	}
}

// TODO TestSurroundedGroupsAreCapturedTogether(t *testing.T) {}

func TestCorneredGroupsAreCapturedTogether(t *testing.T) {
	stones := []game.Stone{
		white(0, 0), white(1, 0), white(2, 0), black(3, 0),
		white(0, 1), black(1, 1), black(2, 1),
	}
	m := game.New(game.BoardSizeTiny, stones...)
	m.Play(black(0, 2))

	// Expect all white stones to be captured
	for _, stone := range stones {
		if stone.Colour != game.White {
			continue
		}
		if got := m.Board().At(stone.X, stone.Y); got != game.None {
			t.Errorf("Expected white stone at position {%d,%d} to be captured, instead got: '%s'",
				stone.X, stone.Y, got)
		}
	}
}

func TestSuicidalMovesAreNotAllowed(t *testing.T) {
	// TODO suicide is allowed in some rulesets
	// https://senseis.xmp.net/?Suicide

	stones := []game.Stone{
		white(2, 1), white(1, 2), white(2, 3), white(3, 2),
	}
	m := game.New(game.BoardSizeTiny, stones...)

	err := m.Play(black(2, 2))
	if expected := game.ErrSuicidalMove; err != expected {
		t.Errorf("Expected: '%v', Got: '%v'", expected, err)
	}
}

func TestCapturesResolveBeforeSuicide(t *testing.T) {
	stones := []game.Stone{
		white(2, 1), white(1, 2), white(2, 3), white(3, 2),
		black(3, 1), black(4, 2), black(3, 3),
	}
	m := game.New(game.BoardSizeTiny, stones...)

	got := m.Play(black(2, 2))
	if got != nil {
		t.Errorf("Expected: '%v', got: '%v'", nil, got)
	}
}

func TestPlayFailsWhenKoIsViolated(t *testing.T) {
	stones := []game.Stone{
		white(2, 1), white(1, 2), white(2, 3), white(3, 2),
		black(3, 1), black(4, 2), black(3, 3),
	}
	m := game.New(game.BoardSizeTiny, stones...)
	m.Play(black(2, 2))

	got := m.Play(white(3, 2))
	if expected := game.ErrViolatesKo; got != expected {
		t.Errorf("Expected: '%v', got: '%v'", expected, got)
	}
}

func black(x, y int) game.Stone {
	return stone(game.Black, x, y)
}

func white(x, y int) game.Stone {
	return stone(game.White, x, y)
}

func stone(c game.Colour, x, y int) game.Stone {
	return game.Stone{Colour: c, Position: game.Position{X: x, Y: y}}
}
