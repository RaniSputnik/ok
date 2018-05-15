package game

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

type group struct {
	Colour    Colour
	Indexes   []int
	Liberties int
}

func (g group) empty() bool {
	return len(g.Indexes) == 0
}
