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

	PreviousPawnFirstMove bool

	// TODO: To later UNDO a Movement, might be necessary to add more parameters (such as Castling abilities)
}

func NewMovement(movingPiece, takingPiece *piece.Piece, from, to position.Position, previousPawnFirstMove bool) Movement {
	return Movement{
		MovingPiece: movingPiece,
		TakingPiece: takingPiece,
		From:        from,
		To:          to,

		PreviousPawnFirstMove: previousPawnFirstMove,
	}
}
