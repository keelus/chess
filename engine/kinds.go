package engine

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
