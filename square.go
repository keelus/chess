package chess

import (
	"errors"
)

type Square struct {
	I, J uint8
}

func newSquare(i, j uint8) Square {
	return Square{
		I: i,
		J: j,
	}
}

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

func (s1 Square) IsEqualTo(s2 Square) bool {
	return s1.I == s2.I && s1.J == s2.J
}

func (s Square) Algebraic() string {
	return string([]rune{rune(s.J) + 'a', '8' - rune(s.I)})
}

func NewSquareFromAlgebraic(algebraic string) (Square, error) {
	if len(algebraic) != 2 {
		return Square{}, errors.New("Invalid algebraic square provided. It must be the format, e.g.: \"d5\"")
	}

	col := algebraic[0]
	row := algebraic[1]

	if col < 'a' || col > 'h' || row < '1' || row > '8' {
		return Square{}, errors.New("Invalid algebraic square provided. It must be the format, e.g.: \"d5\"")
	}

	finalCol := uint8(col - 'a')
	finalRow := uint8('8' - row)

	return newSquare(finalRow, finalCol), nil
}
