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

// Board contains the board state.
type Board struct {
	Size   int
	Stones []Colour // TODO private?
}

// Contains returns whether or not a given x,y position
// is within the boards boundaries.
func (b Board) Contains(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.Size && y < b.Size
}

// At returns the stone colour at the given x,y position.
//
// Returns `None` if the x,y position is empty or
// is outside the board boundaries.
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

	// TODO investigate simply reading the group
	// to see if we've already visited a cell with
	// the given index
	visited := map[int]bool{}

	var walkGroup func(g group, x, y int) group
	walkGroup = func(g group, x, y int) group {
		i := b.index(x, y)
		if visited[i] {
			return g
		}
		visited[i] = true
		g.Indexes = append(g.Indexes, i)

		for _, n := range b.neighbours(x, y) {
			if !b.Contains(n.x, n.y) {
				continue
			}
			switch nColour := b.At(n.x, n.y); nColour {
			case None:
				g.Liberties++
			case g.Colour:
				g = walkGroup(g, n.x, n.y)
			default:
				// Opponent stone, no liberties here
			}
		}
		return g
	}

	g = walkGroup(g, x, y)

	return g
}
