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

func TestPlayAddsAStone(t *testing.T) {
	test := game.New(game.BoardSizeTiny)
	playedStone := game.Stone{game.Black, game.Position{0, 0}}
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
	_, err := test.Play(game.Stone{game.White, game.Position{0, 0}})
	if err != game.ErrNotYourTurn {
		t.Errorf("Expected: '%v', got: '%v'", game.ErrNotYourTurn, err)
	}
}

func TestPlayChangesNext(t *testing.T) {
	m := game.New(game.BoardSizeTiny)
	m, _ = m.Play(game.Stone{game.Black, game.Position{0, 0}})
	if got := m.Next(); got != game.White {
		t.Errorf("Expected: '%s' to go second, got: '%s'", game.White, got)
	}
	m, _ = m.Play(game.Stone{game.White, game.Position{1, 0}})
	if got := m.Next(); got != game.Black {
		t.Errorf("Expected: '%s' to go third, got: '%s'", game.Black, got)
	}
}
