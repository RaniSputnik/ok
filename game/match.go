package game

type Match struct {
	current Board
}

func New(size int, stones ...Stone) Match {
	return Match{
		current: Board{
			Size: size,
		},
	}
}

func (m Match) Board() Board {
	return m.current
}
