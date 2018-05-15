package game

type Move interface {
	Player() Colour
}

type skip struct {
	c Colour
}

func (s skip) Player() Colour {
	return s.c
}
