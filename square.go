package chess

import (
	"errors"
)

// Square represents a coordinate/point in the board.
type Square struct {
	I, J uint8
}

func newSquare(i, j uint8) Square {
	return Square{
		I: i,
		J: j,
	}
}

// NewSquare returns a new instance of Square and an error, if the
// square coordinates are not valid.
//
// Coordinates should be in range of [0, 8)
func NewSquare(i, j uint8) (Square, error) {
	if i > 7 || j > 7 {
		return Square{}, errors.New("Error creating square.")
	}

	return newSquare(i, j), nil
}

func (s *Square) clone() Square {
	return Square{
		s.I,
		s.J,
	}
}

// IsEqualTo reports whether a square is equal to another square.
func (s1 Square) IsEqualTo(s2 Square) bool {
	return s1.I == s2.I && s1.J == s2.J
}

// Algebraic returns the algebraic representation of the board.
//
// It assumes that the square is valid, and inside the board.
func (s Square) Algebraic() string {
	return string([]rune{rune(s.J) + 'a', '8' - rune(s.I)})
}

// NewSquareFromAlgebraic returns a new square from the Pure algebraic
// notation passed. It will return an error if the position is not valid.
//
// Examples:
//
//	NewSquareFromAlgebraic("aa") // returns Square{...}, error
//	NewSquareFromAlgebraic("d2") // returns Square{I:6, J:3}, nil
//	NewSquareFromAlgebraic("h8") // returns Square{I:0, J:7}, nil
func NewSquareFromAlgebraic(algebraic string) (Square, error) {
	if len(algebraic) != 2 {
		return Square{}, errors.New("Invalid algebraic square provided. It must be in the algebraic notation (example: \"d5\")")
	}

	col := algebraic[0]
	row := algebraic[1]

	if col < 'a' || col > 'h' || row < '1' || row > '8' {
		return Square{}, errors.New("Invalid algebraic square provided. It must be in the algebraic notation (example: \"d5\")")
	}

	finalCol := uint8(col - 'a')
	finalRow := uint8('8' - row)

	return newSquare(finalRow, finalCol), nil
}
