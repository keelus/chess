package piece

import (
	"chess/position"
)

type Piece struct {
	Color Color
	Kind  Kind

	IsPawnFirstMovement bool

	Position position.Position
}

func NewPiece(color Color, kind Kind, pos position.Position) Piece {
	return Piece{
		Color: color,
		Kind:  kind,

		IsPawnFirstMovement: kind == Kind_Pawn,

		Position: pos,
	}
}
