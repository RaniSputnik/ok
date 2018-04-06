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
