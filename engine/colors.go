package engine

const COLOR_AMOUNT = 2

type Color uint

const (
	Color_White Color = iota
	Color_Black
	Color_None
)

func (c Color) ToRune() rune {
	if c == Color_White {
		return 'w'
	}
	return 'b'
}

func ColorFromRune(color rune) Color {
	if color == 'w' {
		return Color_White
	}
	return Color_Black
}
