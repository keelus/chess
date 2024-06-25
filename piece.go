package chess

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
