package chess

import "unicode"

// The total amount of piece kind/types available (excluding Kind_None)
const KIND_AMOUNT = 6

// Kind represents a Piece's kind or type
type Kind uint8

const (
	Kind_None   Kind = iota // Not initialized Kind
	Kind_King               // King
	Kind_Queen              // Queen
	Kind_Rook               // Rook
	Kind_Bishop             // Bishop
	Kind_Knight             // Knight
	Kind_Pawn               // Pawn
)

// Piece offsets, directions and other constants

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

// ToRune returns the rune of the piece kind, in lowercase.
//
// Examples:
//   Kind_None.ToRune()  // returns '_'
//   Kind_Queen.ToRune() // returns 'q'
//   Kind_Knight.ToRune() // returns 'n'
func (k Kind) ToRune() rune {
	return []rune{'_', 'k', 'q', 'r', 'b', 'n', 'p'}[k]
}

// ToRuneWith returns the rune of the piece kind, with proper case
// depending on the passed color.
//
// Examples:
//   Kind_None.ToRune(Color_None)  // returns '_'
//   Kind_Queen.ToRune(Color_White) // returns 'Q'
//   Kind_Queen.ToRune(Color_Black) // returns 'q'
func (k Kind) ToRuneWithColor(color Color) rune {
	if k == Kind_None || color == Color_None {
		return '_'
	}

	r := []rune{'k', 'q', 'r', 'b', 'n', 'p'}[k+1]

	if color == Color_White {
		return unicode.ToUpper(r)
	}

	return r
}

// KindFromRune returns the piece Kind of the provided rune.
// The rune can be lowercase or uppercase.
//
// Examples:
//   KindFromRune('Q') // returns Kind_Queen
//   KidnFromRune('q') // returns Kind_Queen
//   KidnFromRune('n') // returns Kind_Knight
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
	case 'p', 'P':
		return Kind_Pawn
	default:
		return Kind_None
	}
}

// KindAndColorFromRune returns the piece Kind and the color
// of the provided rune, based on the character and case.
//
// Examples:
//   KindAndColorFromRune('Q') // returns Kind_Queen, Color_White
//   KindAndColorFromRune('k') // returns Kind_King, Color_Black
func KindAndColorFromRune(kind rune) (Kind, Color) {
	if kind == '_' {
		return Kind_None, Color_None
	}

	var parsedColor Color

	if unicode.IsUpper(kind) {
		parsedColor = Color_White
	} else {
		parsedColor = Color_Black
	}

	return KindFromRune(kind), parsedColor
}
