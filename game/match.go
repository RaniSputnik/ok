package game

type Match struct {
	current Board
}

func New(size int, stones ...Stone) Match {
	return Match{
		current: Board{
			Size:   size,
			Stones: make([]Colour, size*size),
		},
	}
}

func (m Match) Board() Board {
	return m.current
}

func (m Match) Play(move Stone) (Match, error) {
	// TODO validate move

	next := Board{
		Size:   m.current.Size,
		Stones: m.current.Stones,
	}

	// Set the colour of the square
	i := next.index(move.X, move.Y)
	next.Stones[i] = move.Colour

	m.current = next
	return m, nil
}
