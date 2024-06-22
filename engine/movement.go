package engine

import "fmt"

type Movement struct {
	MovingPiece   Piece
	TakingPiece   Piece // Optional
	IsTakingPiece bool

	From Position
	To   Position

	// Next variables refer to the state before this movement have been done
	PawnIsDoublePositionMovement *bool

	IsKingSideCastling  *bool
	IsQueenSideCastling *bool

	CanWhiteQueenSideCastling bool
	CanWhiteKingSideCastling  bool
	CanBlackQueenSideCastling bool
	CanBlackKingSideCastling  bool

	// CanKingSideCastling  bool // TODO: These two only in castle or King moves
	// CanQueenSideCastling bool

	EnPassant *Position

	// TODO: To later UNDO a Movement, might be necessary to add more parameters (such as Castling abilities)
}

func (m *Movement) WithTakingPiece(piece Piece) *Movement {
	m.TakingPiece = piece
	m.IsTakingPiece = true
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

func NewMovement(movingPiece Piece, from, to Position, enPassant *Position, canWhiteQueenSideCastling, canWhiteKingSideCastling, canBlackQueenSideCastling, canBlackKingSideCastling bool) *Movement {
	return &Movement{
		MovingPiece:   movingPiece,
		IsTakingPiece: false,
		From:          from,
		To:            to,
		EnPassant:     enPassant,

		CanWhiteQueenSideCastling: canWhiteQueenSideCastling,
		CanWhiteKingSideCastling:  canWhiteKingSideCastling,
		CanBlackQueenSideCastling: canBlackQueenSideCastling,
		CanBlackKingSideCastling:  canBlackKingSideCastling,
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

	return fmt.Sprintf("%s%s", from.ToAlgebraic(), to.ToAlgebraic())
}
