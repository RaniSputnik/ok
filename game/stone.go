package game

type Stone struct {
	Colour Colour
	Position
}

type Position struct {
	X, Y int
}
