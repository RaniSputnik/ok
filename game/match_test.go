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
	test := game.New(game.BoardSizeTiny)
	playedStone := black(0, 0)
	got, err := test.Play(playedStone)

	if err != nil {
		t.Errorf("Expected err: '<nil>', got: '%v'", err)
	}

	gotStone := got.Board().Stones[0]
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
	test := game.New(game.BoardSizeTiny)
	_, err := test.Play(white(0, 0))
	if err != game.ErrNotYourTurn {
		t.Errorf("Expected: '%v', got: '%v'", game.ErrNotYourTurn, err)
	}
}

func TestPlayChangesNext(t *testing.T) {
	m := game.New(game.BoardSizeTiny)
	m, _ = m.Play(black(0, 0))
	if got := m.Next(); got != game.White {
		t.Errorf("Expected: '%s' to go second, got: '%s'", game.White, got)
	}
	m, _ = m.Play(white(1, 0))
	if got := m.Next(); got != game.Black {
		t.Errorf("Expected: '%s' to go third, got: '%s'", game.Black, got)
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
