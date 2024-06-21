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

func (m *Movement) WithTakingPiece(piece *Piece) *Movement {
	m.TakingPiece = piece
	return m
}

func (m *Movement) WithPawn(isFirstMove bool) *Movement {
	m.PawnIsFirstMove = &isFirstMove
	return m
}

func (m *Movement) WithCastling(isQueenSideMove, isKingSideMove bool) *Movement {
	m.IsKingSideCastling = &isKingSideMove
	m.IsQueenSideCastling = &isQueenSideMove
	return m
}

func NewMovement(movingPiece *Piece, from, to Position, canQueenSideCastling, canKingSideCastling bool) *Movement {
	return &Movement{
		MovingPiece: movingPiece,
		From:        from,
		To:          to,

		CanQueenSideCastling: canQueenSideCastling,
		CanKingSideCastling:  canKingSideCastling,
	}
}
