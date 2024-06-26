package chess

// The total amount of colors available (excluding Color_None)
const COLOR_AMOUNT = 2

// Color represents a Piece's and the side/player to move's color.
type Color uint8

const (
	Color_None  Color = iota // Used to represent an empty "Piece"
	Color_White              // White color
	Color_Black              // Black color
)

// Rune returns the rune of the color.
//
// Examples:
//   Color_None.Rune()  // returns '_'
//   Color_White.Rune() // returns 'w'
func (c Color) Rune() rune {
	return [3]rune{'_', 'w', 'b'}[c]
}

// String returns the string/word of the color.
//
// Examples:
//   Color_None.String()  // returns 'none'
//   Color_White.String() // returns 'white'
func (c Color) String() string {
	return [3]string{"none", "white", "black"}[c]
}

// ColorFromRune returns the Color type of the rune passed.
//
// Examples:
//   ColorFromRune('_') // returns Color_None
//   ColorFromRune('w') // returns Color_White
//   ColorFromRune('b') // returns Color_Black
func ColorFromRune(color rune) Color {
	if color == 'w' {
		return Color_White
	} else if color == 'b' {
		return Color_Black
	}
	return Color_None
}

// Opposite returns the opposite/opponent's color.
//
// Examples:
//   Color_None.Opposite() // returns Color_None
//   Color_White.Opposite() // returns Color_Black
//   Color_Black.Opposite() // returns Color_White
func (c Color) Opposite() Color {
	if c == Color_White {
		return Color_Black
	} else if c == Color_Black {
		return Color_White
	}
	return Color_None
}
