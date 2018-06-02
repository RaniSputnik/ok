package game

import (
	"strings"
)

type Match struct {
	next    Colour
	prev    Board
	current Board

	moves []Move
}

// New creates a new match with the given board
// size and initial state.
func New(size int, stones ...Stone) *Match {
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

	return &m
}

// Next returns the player who has the current turn.
func (m *Match) Next() Colour {
	return m.next
}

// Board returns the current board state.
func (m *Match) Board() Board {
	return m.current
}

// Moves returns the moves that have been played in this game
func (m *Match) Moves() []Move {
	return m.moves
}

// Play applies a move to the current match.
// Returns an error if the move is invalid.
//
// Use m.Play(game.Stone{}) to play a new
// stone on the board.
//
// Use m.Play(game.Skip(colour Colour)) to
// skip the current players turn.
func (m *Match) Play(move Move) error {
	if move.Player() != m.Next() {
		return ErrNotYourTurn
	}

	var err error
	switch v := move.(type) {
	case Stone:
		err = m.playStone(v)
	case skip:
		err = m.skip()

	// If we don't recognise the move, just
	// add it to the list. This should allow
	// us to easily extend with admin events etc.
	default:
		m.moves = append(m.moves, move)
	}

	if err == nil {
		m.moves = append(m.moves, move)
	}
	return err
}

func (m *Match) playStone(move Stone) error {
	if !m.current.Contains(move.X, move.Y) {
		return ErrOutsideBoard
	}
	if m.current.At(move.X, move.Y) != None {
		return ErrPositionOccupied
	}
	nextBoard := m.current.copy()

	// Set the colour of the square
	i := nextBoard.index(move.X, move.Y)
	nextBoard.Stones[i] = move.Colour

	for _, g := range nextBoard.neighbourGroups(move.X, move.Y) {
		if g.Liberties == 0 {
			nextBoard = nextBoard.removeGroup(g)
		}
	}

	g := nextBoard.findGroup(move.X, move.Y)
	if g.Liberties == 0 {
		return ErrSuicidalMove
	}

	if nextBoard.equals(m.prev) {
		return ErrViolatesKo
	}

	m.next = move.Colour.Opponent()
	m.prev = m.current
	m.current = nextBoard
	return nil
}

func (m *Match) skip() error {
	m.next = m.next.Opponent()
	return nil
}

func (m *Match) String() string {
	b := m.Board()
	border := "@" + strings.Repeat("---", b.Size) + "@\n"

	str := "\n" + border
	i := 0
	for y := 0; y < b.Size; y++ {
		str += "|"
		for x := 0; x < b.Size; x++ {

			switch b.Stones[i] {
			case Black:
				str += " X "
			case White:
				str += " O "
			default:
				str += " . "
			}

			i++
		}
		str += "|\n"
	}
	str += border
	return str
}
