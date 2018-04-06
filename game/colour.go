package game

import (
	"fmt"
)

// Colour represents a stone colour.
// Possible values are Black, White and None.
type Colour byte

const (
	None = Colour(iota)
	Black
	White
)

// String returns the Colour name
func (c Colour) String() string {
	switch c {
	case None:
		return "None"
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return fmt.Sprintf("Colour(%d)", c)
	}
}

// Opponent will return the opposite colour
// of the given colour.
//
// Black.Opponent() will return White
// White.Opponent() will return Black
// All other colours don't have a valid
// opponent and will return None.
func (c Colour) Opponent() Colour {
	switch c {
	case Black:
		return White
	case White:
		return Black
	default:
		return None
	}
}
