package engine

type Piece struct {
	Color Color
	Kind  Kind

	IsPawnFirstMovement bool

	Position Position
}

func NewPiece(color Color, kind Kind, pos Position) Piece {
	return Piece{
		Color: color,
		Kind:  kind,

		IsPawnFirstMovement: kind == Kind_Pawn,

		Position: pos,
	}
}

func (p Piece) DeepCopy() Piece {
	newPiece := NewPiece(p.Color, p.Kind, NewPosition(p.Position.I, p.Position.J))
	newPiece.IsPawnFirstMovement = p.IsPawnFirstMovement
	return newPiece
}
