package engine

type Movement struct {
	MovingPiece *Piece
	TakingPiece *Piece // Optional

	MovingPieceCopy Piece
	TakingPieceCopy *Piece

	From Position
	To   Position

	// Next variables refer to the state before this movement have been done
	PawnIsDoublePositionMovement *bool

	IsKingSideCastling  *bool
	IsQueenSideCastling *bool

	CanKingSideCastling  bool // TODO: These two only in castle or King moves
	CanQueenSideCastling bool

	EnPassant *Position

	// TODO: To later UNDO a Movement, might be necessary to add more parameters (such as Castling abilities)
}

func (m *Movement) WithTakingPiece(piece *Piece) *Movement {
	if piece == nil {
		return m
	}

	m.TakingPiece = piece
	takingPieceCopy := piece.DeepCopy()
	m.TakingPieceCopy = &takingPieceCopy
	return m
}

func (m *Movement) WithPawn(isDoublePositionMovement bool) *Movement {
	m.PawnIsDoublePositionMovement = &isDoublePositionMovement
	return m
}

func (m *Movement) WithCastling(isQueenSideMove, isKingSideMove bool) *Movement {
	m.IsKingSideCastling = &isKingSideMove
	m.IsQueenSideCastling = &isQueenSideMove
	return m
}

func NewMovement(movingPiece *Piece, from, to Position, enPassant *Position, canQueenSideCastling, canKingSideCastling bool) *Movement {
	return &Movement{
		MovingPiece:     movingPiece,
		MovingPieceCopy: movingPiece.DeepCopy(),
		From:            from,
		To:              to,
		EnPassant:       enPassant,

		CanQueenSideCastling: canQueenSideCastling,
		CanKingSideCastling:  canKingSideCastling,
	}
}
