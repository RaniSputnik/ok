package game

type Match struct {
	next    Colour
	current Board
}

func New(size int, stones ...Stone) Match {
	// TODO what to do about invalid stone positions?
	// eg. overlapping stones or stones with no liberties.
	m := Match{
		next: Black,
		current: Board{
			Size:   size,
			Stones: make([]Colour, size*size),
		},
	}

	for _, stone := range stones {
		i := m.current.index(stone.X, stone.Y)
		m.current.Stones[i] = stone.Colour
	}

	return m
}

// Next returns the player who has the current turn.
func (m Match) Next() Colour {
	return m.next
}

// Board returns the current board state.
func (m Match) Board() Board {
	return m.current
}

// Play adds a stone to the board, the mutated match
// is returned. Returns an error if the move is invalid.
func (m Match) Play(move Stone) (Match, error) {
	if move.Colour != m.Next() {
		return m, ErrNotYourTurn
	}

	// TODO validate move

	nextBoard := Board{
		Size:   m.current.Size,
		Stones: m.current.Stones,
	}

	// Set the colour of the square
	i := nextBoard.index(move.X, move.Y)
	nextBoard.Stones[i] = move.Colour

	m.next = move.Colour.Opponent()
	m.current = nextBoard
	return m, nil
}
