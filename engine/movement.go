package engine

import "fmt"

type Movement struct {
	MovingPiece   Piece
	TakingPiece   Piece // Optional
	IsTakingPiece bool

	From Point
	To   Point

	// Next variables refer to the state before this movement have been done
	PawnIsDoublePointMovement           *bool
	PawnIsAttackingButNotTakingDiagonal *bool
	PawnPromotionTo                     *Kind

	IsKingSideCastling  *bool
	IsQueenSideCastling *bool
}

func (m *Movement) WithPawnPromotion(newKind Kind) *Movement {
	m.PawnPromotionTo = &newKind
	return m
}

func (m *Movement) WithTakingPiece(piece Piece) *Movement {
	m.TakingPiece = piece
	m.IsTakingPiece = true
	return m
}

func (m *Movement) WithPawn(isDoublePointMovement, attackingButNotTakingDiagonal bool) *Movement {
	m.PawnIsDoublePointMovement = &isDoublePointMovement
	m.PawnIsAttackingButNotTakingDiagonal = &attackingButNotTakingDiagonal
	return m
}

func (m *Movement) WithCastling(isQueenSideMove, isKingSideMove bool) *Movement {
	m.IsKingSideCastling = &isKingSideMove
	m.IsQueenSideCastling = &isQueenSideMove
	return m
}

func NewMovement(movingPiece Piece, from, to Point /*, enPassant *Point, canWhiteQueenSideCastling, canWhiteKingSideCastling, canBlackQueenSideCastling, canBlackKingSideCastling bool*/) *Movement {
	return &Movement{
		MovingPiece:   movingPiece,
		IsTakingPiece: false,
		From:          from,
		To:            to,
		// EnPassant:     enPassant,

		// CanWhiteQueenSideCastling: canWhiteQueenSideCastling,
		// CanWhiteKingSideCastling:  canWhiteKingSideCastling,
		// CanBlackQueenSideCastling: canBlackQueenSideCastling,
		// CanBlackKingSideCastling:  canBlackKingSideCastling,
		// CanQueenSideCastling: canQueenSideCastling,
		// CanKingSideCastling:  canKingSideCastling,
	}
}

func (m Movement) ToString() string {
	return fmt.Sprintf("Piece [color: %c, kind: %c] moves from (%d, %d) to (%d, %d) [takes: %t].", m.MovingPiece.Color.ToRune(), m.MovingPiece.Kind.ToRune(), m.From.I, m.From.J, m.To.I, m.To.J, m.IsTakingPiece)
}

func (m Movement) ToAlgebraic() string {
	from := m.From
	to := m.To

	algebraic := fmt.Sprintf("%s%s", from.ToAlgebraic(), to.ToAlgebraic())
	if m.PawnPromotionTo != nil {
		algebraic = fmt.Sprintf("%s%c", algebraic, (*m.PawnPromotionTo).ToRune())
	}

	return algebraic
}
