package piece

type Piece struct {
	Color Color
	Kind  Kind
}

func NewPiece(color Color, kind Kind) Piece {
	return Piece{
		Color: color,
		Kind:  kind,
	}
}
