package game

const (
	// BoardSizeTiny is the recommended
	// board size for new players.
	BoardSizeTiny = 9
	// BoardSizeSmall is for a quicker game.
	BoardSizeSmall = 13
	// BoardSizeNormal is the standard 19x19 board size.
	BoardSizeNormal = 19
)

type Board struct {
	Size   int
	Stones []Colour
}

func (b Board) Contains(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.Size && y < b.Size
}

func (b Board) At(x, y int) Colour {
	if !b.Contains(x, y) {
		return None
	}
	i := b.index(x, y)
	return b.Stones[i]
}

func (b Board) index(x, y int) int {
	return x + y*b.Size
}

func (b Board) pos(index int) (x, y int) {
	return x % b.Size, y / b.Size
}
