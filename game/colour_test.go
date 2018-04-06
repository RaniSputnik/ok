package game_test

import (
	"testing"

	"github.com/RaniSputnik/ko2/game"
)

func TestColourString(t *testing.T) {
	testCases := []struct {
		Colour game.Colour
		Expect string
	}{
		{game.None, "None"},
		{game.Black, "Black"},
		{game.White, "White"},
		{game.Colour(5), "Colour(5)"},
		{game.Colour(7), "Colour(7)"},
	}

	for _, test := range testCases {
		got := test.Colour.String()
		if got != test.Expect {
			t.Errorf("Expected: '%s', Got: '%s'", test.Expect, got)
		}
	}
}

func TestColourOpponent(t *testing.T) {
	testCases := []struct {
		Colour game.Colour
		Expect game.Colour
	}{
		{game.None, game.None},
		{game.Black, game.White},
		{game.White, game.Black},
		{game.Colour(5), game.None},
		{game.Colour(7), game.None},
	}

	for _, test := range testCases {
		got := test.Colour.Opponent()
		if got != test.Expect {
			t.Errorf("Expected: '%s', Got: '%s'", test.Expect, got)
		}
	}
}
