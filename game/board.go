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
