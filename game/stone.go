package game

type Stone struct {
	Colour Colour
	Position
}

type Position struct {
	X, Y int
}

type group struct {
	Colour    Colour
	Indexes   []int
	Liberties int
}

func (g group) empty() bool {
	return len(g.Indexes) == 0
}
