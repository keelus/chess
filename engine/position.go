package engine

import "fmt"

type Position struct {
	I int
	J int
}

func NewPosition(i, j int) Position {
	return Position{
		I: i,
		J: j,
	}
}

func (p Position) ToAlgebraic() string {
	row := p.I
	col := p.J

	cols := [8]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	rows := [8]rune{'8', '7', '6', '5', '4', '3', '2', '1'}

	return fmt.Sprintf("%c%c", cols[col], rows[row])
}
