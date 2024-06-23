package engine

type Piece struct {
	Color Color
	Kind  Kind

	IsPawnFirstMovement bool

	Point Point
}

func NewPiece(color Color, kind Kind, pos Point) Piece {
	return Piece{
		Color: color,
		Kind:  kind,

		IsPawnFirstMovement: kind == Kind_Pawn,

		Point: pos,
	}
}

func (p Piece) DeepCopy() Piece {
	newPiece := NewPiece(p.Color, p.Kind, NewPoint(p.Point.I, p.Point.J))
	newPiece.IsPawnFirstMovement = p.IsPawnFirstMovement
	return newPiece
}
