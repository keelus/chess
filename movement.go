package chess

import (
	"errors"
	"strings"
)

// Movement represents a piece movement to do in a chess position.
type Movement struct {
	movingPiece   Piece
	takingPiece   Piece // Optional
	isTakingPiece bool

	fromSq Square
	toSq   Square

	isDoublePawnPush bool
	pawnPromotionTo  *Kind

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

func (m *Movement) withPawn(isDoublePawnPush bool) *Movement {
	m.isDoublePawnPush = isDoublePawnPush
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

// Algebraic returns the Pure algebraic notation of the movement, as string.
//
// Examples outputs:
//
//	"d2d3"
//	"f7f8q"
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

// FromSquare returns the initial Square of the movement.
func (m Movement) FromSquare() Square {
	return m.fromSq
}

// ToSquare returns the final Square of the movement.
func (m Movement) ToSquare() Square {
	return m.toSq
}

// MovingPiece returns a copy of the Piece that is making the movement.
//
// In case of castling, it would return the King. You can check for castling
// using the movement.IsCastling() function.
func (m Movement) MovingPiece() Piece {
	return m.movingPiece
}

// IsTakingPiece reports whether this movement takes a piece or not.
func (m Movement) IsTakingPiece() bool {
	return m.isTakingPiece
}

// TakingPiece returns the taking/capturing piece in this movement.
//
// If the movement does not take a piece, it will return an empty Piece and the error.
func (m Movement) TakingPiece() (Piece, error) {
	if !m.isTakingPiece {
		return Piece{}, errors.New("This movement is not taking any piece.")
	}
	return m.takingPiece, nil
}

// IsDoublePawnPush reports whether the movement is a pawn's double push or not.
func (m Movement) IsDoublePawnPush() bool {
	return m.isDoublePawnPush
}

// IsPawnPromotion reports whether the movement promotes a pawn or not.
func (m Movement) IsPawnPromotion() bool {
	return m.pawnPromotionTo != nil
}

// IsPawnPromotion returns the new Kind the pawn is being promoted to.
//
// If the movement is not a promotion, it will return Kind_None and the error.
func (m Movement) PawnPromotion() (Kind, error) {
	if m.pawnPromotionTo == nil {
		return Kind_None, errors.New("This movement does not promote a pawn.")
	}
	return *m.pawnPromotionTo, nil
}

// IsQueenSideCastling reports whether the movement is a queenside castling or not.
func (m Movement) IsQueenSideCastling() bool {
	return m.isQueenSideCastling
}

// IsKingSideCastling reports whether the movement is a kingside castling or not.
func (m Movement) IsKingSideCastling() bool {
	return m.isKingSideCastling
}
