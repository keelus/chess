package engine

type Position struct {
	Board [8][8]Piece

	PlayerToMove Color

	CanKingCastling  map[Color]bool
	CanQueenCastling map[Color]bool

	EnPassant *Point
}
