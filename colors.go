package chess

const COLOR_AMOUNT = 2

type Color uint8

const (
	Color_None Color = iota
	Color_White
	Color_Black
)

func (c Color) ToRune() rune {
	return []rune{'_', 'w', 'b'}[c]
}

func ColorFromRune(color rune) Color {
	if color == 'w' {
		return Color_White
	} else if color == 'b' {
		return Color_Black
	}
	return '_'
}

func (c Color) Opposite() Color {
	if c == Color_White {
		return Color_Black
	} else if c == Color_Black {
		return Color_White
	}
	return Color_None
}
