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

// Rune returns the rune of the piece.
//
// Examples:
//   A black pawn // returns 'p'
//   A white king // returns 'K'
func (p Piece) Rune() rune {
	return p.Kind.RuneWithColor(p.Color)
}

// Unicode returns the unicode rune of the piece.
//
// Examples:
//   A black pawn // returns '♟'
//   A white king // returns '♕'
func (p Piece) Unicode() rune {
	return p.Kind.UnicodeWithColor(p.Color)
}

// String returns the color and name of the piece.
//
// Examples:
//   "white knight"
//   "black bishop"
func (p Piece) String() string {
	return p.Color.String() + " " + p.Kind.String()
}
