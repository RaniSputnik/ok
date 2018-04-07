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

func (b Board) removeGroup(g group) Board {
	for _, pos := range g.Indexes {
		b.Stones[pos] = None
	}
	return b
}

func (b Board) neighbours(x, y int) [4]struct{ x, y int } {
	return [4]struct{ x, y int }{
		{x - 1, y}, {x, y - 1}, {x + 1, y}, {x, y + 1},
	}
}

func (b Board) neighbourGroups(x, y int) []group {
	groups := make([]group, 0, 4)
	for _, n := range b.neighbours(x, y) {
		g := b.findGroup(n.x, n.y)
		if !g.empty() {
			// TODO maybe de-dupe
			groups = append(groups, g)
		}
	}

	return groups
}

func (b Board) findGroup(x, y int) group {
	g := group{Colour: b.At(x, y)}
	// Either outside board or there is
	// no stone at this position
	if g.Colour == None {
		return g
	}

	i := b.index(x, y)
	g.Indexes = []int{i}

	// TODO walk all neighbours of the same colour

	liberties := 0
	for _, n := range b.neighbours(x, y) {
		if !b.Contains(n.x, n.y) {
			continue
		}
		if nColour := b.At(n.x, n.y); nColour == None {
			liberties++
		}
	}
	g.Liberties = liberties

	return g
}
