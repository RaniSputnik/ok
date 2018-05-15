package game

type Move interface {
	Player() Colour
}

type Stone struct {
	Colour Colour
	Position
}

func (s Stone) Player() Colour {
	return s.Colour
}

type Position struct {
	X, Y int
}

// Skip returns a move that will skip
// the given players go.
func Skip(player Colour) Move {
	return skip{player}
}

type skip struct {
	c Colour
}

func (s skip) Player() Colour {
	return s.c
}
