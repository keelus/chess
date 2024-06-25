package chess

import "unicode"

const KIND_AMOUNT = 6

type Kind uint8

const (
	Kind_None Kind = iota
	Kind_King
	Kind_Queen
	Kind_Rook
	Kind_Bishop
	Kind_Knight
	Kind_Pawn
)

var kingOffsets = [8][2]int8{{-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}}
var pawnMoveRowDirections = map[Color]int8{
	Color_White: -1,
	Color_Black: 1,
}
var pawnAttackOffsets = map[Color][2][2]int8{
	Color_White: {{-1, -1}, {-1, 1}},
	Color_Black: {{1, -1}, {1, 1}},
}
var knightOffsets = [8][2]int8{{-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}}
var bishopDirections = [4][2]int8{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}}
var rookDirections = [4][2]int8{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

var pawnPromotionRows = map[Color]uint8{
	Color_White: 0,
	Color_Black: 7,
}
var pawnStartingRows = map[Color]uint8{
	Color_White: 6,
	Color_Black: 1,
}

var promotableKinds = [4]Kind{Kind_Queen, Kind_Rook, Kind_Bishop, Kind_Knight}

func (t Kind) ToRune() rune {
	return []rune{'_', 'k', 'q', 'r', 'b', 'n', 'p'}[t]
}

func KindFromRune(kind rune) Kind {
	switch kind {
	case 'k', 'K':
		return Kind_King
	case 'q', 'Q':
		return Kind_Queen
	case 'r', 'R':
		return Kind_Rook
	case 'b', 'B':
		return Kind_Bishop
	case 'n', 'N':
		return Kind_Knight
	default:
		return Kind_Pawn
	}
}

func KindAndColorFromRune(kind rune) (Kind, Color) {
	var parsedKind Kind
	var parsedColor Color

	if unicode.IsUpper(kind) {
		parsedColor = Color_White
	} else {
		parsedColor = Color_Black
	}

	parsedKind = KindFromRune(kind)

	return parsedKind, parsedColor
}
