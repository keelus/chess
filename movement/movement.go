package movement

import (
	"chess/piece"
	"chess/position"
)

type Movement struct {
	MovingPiece *piece.Piece
	TakingPiece *piece.Piece // Optional

	From position.Position
	To   position.Position

	// TODO: To later UNDO a Movement, might be necessary to add more parameters (such as Castling abilities)
}

func NewMovement() Movement {
	panic("TODO")
	return Movement{}
}
