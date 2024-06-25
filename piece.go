package chess

// Piece represents a chess Piece. It can be initialized or not.
//
// If not initialized, both Color and Kind will be of type
// none (Color_None, Kind_None)
type Piece struct {
	Color Color
	Kind  Kind

	Square Square
}

func newPiece(color Color, kind Kind, pos Square) Piece {
	return Piece{
		Color:  color,
		Kind:   kind,
		Square: pos,
	}
}

func (p Piece) clone() Piece {
	return newPiece(p.Color, p.Kind, newSquare(p.Square.I, p.Square.J))
}
