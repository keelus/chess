package chess

import (
	"fmt"
	"strings"
)

type Movement struct {
	movingPiece   Piece
	takingPiece   Piece // Optional
	isTakingPiece bool

	fromSq Square
	toSq   Square

	pawnIsDoubleSquareMovement bool
	pawnPromotionTo            *Kind

	isQueenSideCastling bool
	isKingSideCastling  bool
}

func (m *Movement) withPawnPromotion(newKind Kind) *Movement {
	m.pawnPromotionTo = &newKind
	return m
}

func (m *Movement) withTakingPiece(piece Piece) *Movement {
	m.takingPiece = piece
	m.isTakingPiece = true
	return m
}

func (m *Movement) withPawn(isDoubleSquareMovement bool) *Movement {
	m.pawnIsDoubleSquareMovement = isDoubleSquareMovement
	return m
}

func (m *Movement) withCastling(isQueenSideMove, isKingSideMove bool) *Movement {
	m.isKingSideCastling = isKingSideMove
	m.isQueenSideCastling = isQueenSideMove
	return m
}

func newMovement(movingPiece Piece, fromSquare, toSquare Square) *Movement {
	return &Movement{
		movingPiece:   movingPiece,
		isTakingPiece: false,
		fromSq:        fromSquare,
		toSq:          toSquare,
	}
}

func (m Movement) debug() string {
	return fmt.Sprintf("Piece [color: %c, kind: %c] moves from (%d, %d) to (%d, %d) [takes: %t].", m.movingPiece.Color.ToRune(), m.movingPiece.Kind.ToRune(), m.fromSq.I, m.fromSq.J, m.toSq.I, m.toSq.J, m.isTakingPiece)
}

func (m Movement) Algebraic() string {
	from := m.fromSq
	to := m.toSq

	var sb strings.Builder

	sb.WriteString(from.Algebraic())
	sb.WriteString(to.Algebraic())

	if m.pawnPromotionTo != nil {
		sb.WriteRune((*m.pawnPromotionTo).ToRune())
	}

	return sb.String()
}

func (m Movement) FromSquare() Square {
	return m.fromSq
}

func (m Movement) ToSquare() Square {
	return m.toSq
}

func (m Movement) IsCapturing() bool {
	return m.isTakingPiece
}
