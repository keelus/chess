package engine

type Movement struct {
	MovingPiece *Piece
	TakingPiece *Piece // Optional

	From Position
	To   Position

	// Next variables refer to the state before this movement have been done
	PawnIsFirstMove *bool

	IsKingSideCastling  *bool
	IsQueenSideCastling *bool

	CanKingSideCastling  bool // TODO: These two only in castle or King moves
	CanQueenSideCastling bool

	// TODO: To later UNDO a Movement, might be necessary to add more parameters (such as Castling abilities)
}

func NewMovement(movingPiece, takingPiece *Piece, from, to Position, canKingSideCastling, canQueenSideCastling bool) Movement {
	return Movement{
		MovingPiece:          movingPiece,
		TakingPiece:          takingPiece,
		From:                 from,
		To:                   to,
		CanKingSideCastling:  canKingSideCastling,
		CanQueenSideCastling: canQueenSideCastling,
	}
}

func NewPawnMovement(movingPiece, takingPiece *Piece, from, to Position, canKingSideCastling, canQueenSideCastling, pawnIsFirstMove bool) Movement {
	newMovement := NewMovement(movingPiece, takingPiece, from, to, canKingSideCastling, canQueenSideCastling)
	newMovement.PawnIsFirstMove = &pawnIsFirstMove
	return newMovement
}

func NewCastlingMovement(movingPiece *Piece, from, to Position, canKingSideCastling, canQueenSideCastling, isKingSideCastling, isQueenSideCastling bool) Movement {
	newMovement := NewMovement(movingPiece, nil, from, to, canKingSideCastling, canQueenSideCastling)
	newMovement.IsKingSideCastling = &isKingSideCastling
	newMovement.IsQueenSideCastling = &isQueenSideCastling
	return newMovement
}
