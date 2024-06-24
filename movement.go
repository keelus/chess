package engine

import (
	"fmt"
	"strings"
)

type Movement struct {
	movingPiece   Piece
	takingPiece   Piece // Optional
	isTakingPiece bool

	from Square
	to   Square

	pawnIsDoubleSquareMovement bool
	pawnPromotionTo            *Kind

	isKingSideCastling  bool
	isQueenSideCastling bool
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

func newMovement(movingPiece Piece, from, to Square /*, enPassant *Square, canWhiteQueenSideCastling, canWhiteKingSideCastling, canBlackQueenSideCastling, canBlackKingSideCastling bool*/) *Movement {
	return &Movement{
		movingPiece:   movingPiece,
		isTakingPiece: false,
		from:          from,
		to:            to,
	}
}

func (m Movement) debug() string {
	return fmt.Sprintf("Piece [color: %c, kind: %c] moves from (%d, %d) to (%d, %d) [takes: %t].", m.movingPiece.Color.ToRune(), m.movingPiece.Kind.ToRune(), m.from.I, m.from.J, m.to.I, m.to.J, m.isTakingPiece)
}

func (m Movement) Algebraic() string {
	from := m.from
	to := m.to

	var sb strings.Builder

	sb.WriteString(from.ToAlgebraic())
	sb.WriteString(to.ToAlgebraic())

	if m.pawnPromotionTo != nil {
		sb.WriteRune((*m.pawnPromotionTo).ToRune())
	}

	return sb.String()
}

func (m Movement) From() Square {
	return m.from
}

func (m Movement) To() Square {
	return m.to
}

func (m Movement) IsCapturing() bool {
	return m.isTakingPiece
}
