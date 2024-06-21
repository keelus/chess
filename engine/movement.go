package engine

type Movement struct {
	MovingPiece *Piece
	TakingPiece *Piece // Optional

	From Position
	To   Position

	PreviousPawnFirstMove bool

	// TODO: To later UNDO a Movement, might be necessary to add more parameters (such as Castling abilities)
}

func NewMovement(movingPiece, takingPiece *Piece, from, to Position, previousPawnFirstMove bool) Movement {
	return Movement{
		MovingPiece: movingPiece,
		TakingPiece: takingPiece,
		From:        from,
		To:          to,

		PreviousPawnFirstMove: previousPawnFirstMove,
	}
}
