package engine

const COLOR_AMOUNT = 2

type Color uint8

const (
	Color_None Color = iota
	Color_White
	Color_Black
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
